import React, { useState } from 'react';
import type { Term } from '../types';
import { stringToColor } from '../utils/colorUtils';
import './TermList.css';

interface TermListProps {
  terms: Term[];
  allTerms?: Term[];
  onTermClick: (term: Term) => void;
  onAddTerm: () => void;
  onSearch: (search: string, category: string) => void;
}

const TermList: React.FC<TermListProps> = ({
  terms,
  allTerms,
  onTermClick,
  onAddTerm,
  onSearch,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');

  const safeTerms = terms || [];
  const categoriesSource = allTerms && allTerms.length > 0 ? allTerms : safeTerms;
  const categories = Array.from(new Set(categoriesSource.map((t) => t.category))).sort();

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchQuery(value);
    onSearch(value, selectedCategory);
  };

  const handleCategoryChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    setSelectedCategory(value);
    onSearch(searchQuery, value);
  };

  const handleClearFilters = () => {
    setSearchQuery('');
    setSelectedCategory('');
    onSearch('', '');
  };

  return (
    <div className="term-list">
      <div className="term-list-header">
        <h2>Glossary Terms</h2>
        <button className="primary" onClick={onAddTerm}>
          + Add Term
        </button>
      </div>

      <div className="term-list-filters">
        <div className="filter-group">
          <input
            type="text"
            placeholder="Search terms..."
            value={searchQuery}
            onChange={handleSearchChange}
            className="search-input"
          />
        </div>
        <div className="filter-group">
          <select
            value={selectedCategory}
            onChange={handleCategoryChange}
            className="category-select"
          >
            <option value="">All Categories</option>
            {categories.map((cat) => (
              <option key={cat} value={cat}>
                {cat}
              </option>
            ))}
          </select>
        </div>
        {(searchQuery || selectedCategory) && (
          <button className="secondary" onClick={handleClearFilters}>
            Clear
          </button>
        )}
      </div>

      <div className="term-list-items">
        {safeTerms.length === 0 ? (
          <div className="no-terms">
            <p>No terms found</p>
          </div>
        ) : (
          safeTerms.map((term) => {
            const categoryColor = stringToColor(term.category);
            return (
              <div
                key={term.id}
                className="term-item"
                onClick={() => onTermClick(term)}
                style={{ borderLeftColor: categoryColor }}
              >
                <div className="term-item-header">
                  <h3>{term.name}</h3>
                  <span
                    className="term-category"
                    style={{
                      backgroundColor: categoryColor,
                      color: '#ffffff',
                    }}
                  >
                    {term.category}
                  </span>
                </div>
                <p className="term-definition">
                  {term.definition.length > 150
                    ? `${term.definition.substring(0, 150)}...`
                    : term.definition}
                </p>
              </div>
            );
          })
        )}
      </div>

      <div className="term-list-footer">
        <p>{safeTerms.length} term(s) found</p>
      </div>
    </div>
  );
};

export default TermList;

