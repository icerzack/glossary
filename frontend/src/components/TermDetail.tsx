import React, { useState, useEffect } from 'react';
import type { Term, CreateTermRequest, UpdateTermRequest } from '../types';
import './TermDetail.css';

interface TermDetailProps {
  term: Term | null;
  isEditing: boolean;
  onClose: () => void;
  onSave: (data: CreateTermRequest | UpdateTermRequest) => void;
  onDelete?: () => void;
}

const TermDetail: React.FC<TermDetailProps> = ({
  term,
  isEditing,
  onClose,
  onSave,
  onDelete,
}) => {
  const [formData, setFormData] = useState<CreateTermRequest>({
    name: '',
    definition: '',
    category: '',
    source: '',
    source_url: '',
  });

  useEffect(() => {
    if (term) {
      setFormData({
        name: term.name,
        definition: term.definition,
        category: term.category,
        source: term.source,
        source_url: term.source_url,
      });
    } else {
      setFormData({
        name: '',
        definition: '',
        category: '',
        source: '',
        source_url: '',
      });
    }
  }, [term]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(formData);
  };

  const isNewTerm = !term;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{isNewTerm ? 'Add New Term' : isEditing ? 'Edit Term' : 'Term Details'}</h2>
          <button className="modal-close" onClick={onClose}>
            ×
          </button>
        </div>

        {isEditing || isNewTerm ? (
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name">Term Name *</label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="definition">Definition *</label>
              <textarea
                id="definition"
                name="definition"
                value={formData.definition}
                onChange={handleChange}
                required
                rows={4}
              />
            </div>

            <div className="form-group">
              <label htmlFor="category">Category *</label>
              <input
                type="text"
                id="category"
                name="category"
                value={formData.category}
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="source">Source</label>
              <input
                type="text"
                id="source"
                name="source"
                value={formData.source}
                onChange={handleChange}
              />
            </div>

            <div className="form-group">
              <label htmlFor="source_url">Source URL</label>
              <input
                type="url"
                id="source_url"
                name="source_url"
                value={formData.source_url}
                onChange={handleChange}
              />
            </div>

            <div className="modal-footer">
              <button type="button" className="secondary" onClick={onClose}>
                Cancel
              </button>
              <button type="submit" className="primary">
                {isNewTerm ? 'Create' : 'Save'}
              </button>
            </div>
          </form>
        ) : (
          <div className="term-detail-view">
            <div className="detail-section">
              <h3>Name</h3>
              <p>{term?.name}</p>
            </div>

            <div className="detail-section">
              <h3>Category</h3>
              <span className="category-badge">{term?.category}</span>
            </div>

            <div className="detail-section">
              <h3>Definition</h3>
              <p>{term?.definition}</p>
            </div>

            {term?.source && (
              <div className="detail-section">
                <h3>Source</h3>
                <p>{term.source}</p>
              </div>
            )}

            {term?.source_url && (
              <div className="detail-section">
                <h3>Source URL</h3>
                <a
                  href={term.source_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="source-link"
                >
                  {term.source_url}
                </a>
              </div>
            )}

            <div className="detail-section">
              <h3>Metadata</h3>
              <p className="metadata">
                Created: {new Date(term?.created_at || '').toLocaleDateString()}
                <br />
                Updated: {new Date(term?.updated_at || '').toLocaleDateString()}
              </p>
            </div>

            <div className="modal-footer">
              {onDelete && (
                <button className="danger" onClick={onDelete}>
                  Delete
                </button>
              )}
              <button className="secondary" onClick={onClose}>
                Close
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default TermDetail;

