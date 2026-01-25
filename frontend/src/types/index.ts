export interface Term {
  id: number;
  name: string;
  definition: string;
  category: string;
  source: string;
  source_url: string;
  created_at: string;
  updated_at: string;
}

export interface Relationship {
  id: number;
  source_term_id: number;
  target_term_id: number;
  type: string;
  description: string;
  created_at: string;
}

export interface GraphNode {
  id: number;
  name: string;
  definition: string;
  category: string;
}

export interface GraphEdge {
  id: number;
  source: number;
  target: number;
  type: string;
  description: string;
}

export interface Graph {
  nodes: GraphNode[];
  edges: GraphEdge[];
}

export interface CreateTermRequest {
  name: string;
  definition: string;
  category: string;
  source: string;
  source_url: string;
}

export interface UpdateTermRequest {
  name: string;
  definition: string;
  category: string;
  source: string;
  source_url: string;
}

export interface CreateRelationshipRequest {
  source_term_id: number;
  target_term_id: number;
  type: string;
  description: string;
}

