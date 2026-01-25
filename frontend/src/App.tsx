import { useState, useEffect, useCallback } from 'react';
import MindMap from './components/MindMap';
import TermList from './components/TermList';
import TermDetail from './components/TermDetail';
import { termsApi, graphApi } from './services/api';
import type { Term, Graph, CreateTermRequest, UpdateTermRequest } from './types';
import './App.css';

function App() {
  const [allTerms, setAllTerms] = useState<Term[]>([]);
  const [terms, setTerms] = useState<Term[]>([]);
  const [graph, setGraph] = useState<Graph>({ nodes: [], edges: [] });
  const [selectedTerm, setSelectedTerm] = useState<Term | null>(null);
  const [showDetail, setShowDetail] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const [termsData, graphData] = await Promise.all([
        termsApi.getAll(),
        graphApi.get(),
      ]);
      setAllTerms(termsData);
      setTerms(termsData);
      setGraph(graphData);
    } catch (err) {
      setError('Failed to load data. Please try again.');
      console.error('Error loading data:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleSearch = async (search: string, category: string) => {
    try {
      const data = await termsApi.getAll(search, category);
      setTerms(data);
    } catch (err) {
      console.error('Error searching terms:', err);
    }
  };

  const handleNodeClick = async (nodeId: number) => {
    try {
      const term = await termsApi.getById(nodeId);
      setSelectedTerm(term);
      setIsEditing(false);
      setShowDetail(true);
    } catch (err) {
      console.error('Error loading term:', err);
    }
  };

  const handleTermClick = (term: Term) => {
    setSelectedTerm(term);
    setIsEditing(false);
    setShowDetail(true);
  };

  const handleAddTerm = () => {
    setSelectedTerm(null);
    setIsEditing(true);
    setShowDetail(true);
  };

  const handleSaveTerm = async (data: CreateTermRequest | UpdateTermRequest) => {
    try {
      if (selectedTerm) {
        await termsApi.update(selectedTerm.id, data as UpdateTermRequest);
      } else {
        await termsApi.create(data as CreateTermRequest);
      }
      setShowDetail(false);
      setSelectedTerm(null);
      setIsEditing(false);
      await loadData();
    } catch (err) {
      console.error('Error saving term:', err);
      alert('Failed to save term. Please try again.');
    }
  };

  const handleDeleteTerm = async () => {
    if (!selectedTerm) return;
    
    if (!window.confirm(`Are you sure you want to delete "${selectedTerm.name}"?`)) {
      return;
    }

    try {
      await termsApi.delete(selectedTerm.id);
      setShowDetail(false);
      setSelectedTerm(null);
      await loadData();
    } catch (err) {
      console.error('Error deleting term:', err);
      alert('Failed to delete term. Please try again.');
    }
  };

  const handleCloseDetail = () => {
    setShowDetail(false);
    setSelectedTerm(null);
    setIsEditing(false);
  };

  if (loading) {
    return (
      <div className="app-loading">
        <div className="spinner"></div>
        <p>Loading glossary...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="app-error">
        <h2>Error</h2>
        <p>{error}</p>
        <button className="primary" onClick={loadData}>
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Глоссарий терминов</h1>
        <p>Кузнецов М. А. P4209</p>
      </header>

      <div className="app-content">
        <aside className="app-sidebar">
          <TermList
            terms={terms}
            allTerms={allTerms}
            onTermClick={handleTermClick}
            onAddTerm={handleAddTerm}
            onSearch={handleSearch}
          />
        </aside>

        <main className="app-main">
          <MindMap graph={graph} onNodeClick={handleNodeClick} />
        </main>
      </div>

      {showDetail && (
        <TermDetail
          term={selectedTerm}
          isEditing={isEditing}
          onClose={handleCloseDetail}
          onSave={handleSaveTerm}
          onDelete={selectedTerm ? handleDeleteTerm : undefined}
        />
      )}
    </div>
  );
}

export default App;

