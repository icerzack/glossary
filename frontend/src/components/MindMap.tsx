import React, { useCallback, useEffect, useMemo, useState } from 'react';
import ReactFlow, {
  Node,
  Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  Connection,
  MarkerType,
  BackgroundVariant,
} from 'reactflow';
import 'reactflow/dist/style.css';
import type { Graph, GraphNode as ApiGraphNode } from '../types';
import { stringToColor } from '../utils/colorUtils';
import { relationshipsApi } from '../services/api';
import './MindMap.css';

interface MindMapProps {
  graph: Graph;
  onNodeClick: (nodeId: number) => void;
  onGraphUpdate?: () => void;
}

const getNodeSize = (name: string): { width: number; height: number } => {
  const baseWidth = 140;
  const baseHeight = 60;
  const charWidth = 8;
  const width = Math.max(baseWidth, name.length * charWidth + 40);
  return { width, height: baseHeight };
};

const calculateLayout = (
  nodes: ApiGraphNode[],
  edges: { source: number; target: number }[]
): Record<number, { x: number; y: number }> => {
  const nodePositions: Record<number, { x: number; y: number }> = {};
  
  if (nodes.length === 0) {
    return nodePositions;
  }

  const incomingEdges: Record<number, number[]> = {};
  const outgoingEdges: Record<number, number[]> = {};
  const nodeIds = new Set<number>();

  nodes.forEach((node) => {
    nodeIds.add(node.id);
    incomingEdges[node.id] = [];
    outgoingEdges[node.id] = [];
  });

  edges.forEach((edge) => {
    if (nodeIds.has(edge.source) && nodeIds.has(edge.target)) {
      outgoingEdges[edge.source].push(edge.target);
      incomingEdges[edge.target].push(edge.source);
    }
  });

  const rootNodes = nodes.filter((node) => incomingEdges[node.id].length === 0);
  
  if (rootNodes.length === 0) {
    const minIncoming = Math.min(...nodes.map((node) => incomingEdges[node.id].length));
    rootNodes.push(...nodes.filter((node) => incomingEdges[node.id].length === minIncoming));
  }

  const levels: number[][] = [];
  const nodeLevel: Record<number, number> = {};
  const visited = new Set<number>();

  levels[0] = rootNodes.map((node) => node.id);
  rootNodes.forEach((node) => {
    nodeLevel[node.id] = 0;
    visited.add(node.id);
  });

  let currentLevel = 0;
  while (levels[currentLevel] && levels[currentLevel].length > 0) {
    const nextLevel: number[] = [];
    
    levels[currentLevel].forEach((nodeId) => {
      outgoingEdges[nodeId].forEach((targetId) => {
        if (!visited.has(targetId)) {
          const allIncomingVisited = incomingEdges[targetId].every((sourceId) =>
            visited.has(sourceId)
          );
          
          if (allIncomingVisited) {
            nextLevel.push(targetId);
            nodeLevel[targetId] = currentLevel + 1;
            visited.add(targetId);
          }
        }
      });
    });

    if (nextLevel.length > 0) {
      levels[currentLevel + 1] = nextLevel;
      currentLevel++;
    } else {
      break;
    }
  }

  nodes.forEach((node) => {
    if (!visited.has(node.id)) {
      const level = currentLevel + 1;
      if (!levels[level]) {
        levels[level] = [];
      }
      levels[level].push(node.id);
      nodeLevel[node.id] = level;
    }
  });

  const levelHeight = 200;
  const nodeSpacing = 250;
  const startX = 400;
  const startY = 100;

  levels.forEach((levelNodes, levelIndex) => {
    const y = startY + levelIndex * levelHeight;
    const totalWidth = (levelNodes.length - 1) * nodeSpacing;
    const startXForLevel = startX - totalWidth / 2;

    const sortedNodes = [...levelNodes].sort((a, b) => {
      const aAvgSourceX =
        incomingEdges[a].length > 0
          ? incomingEdges[a].reduce((sum, sourceId) => {
              const pos = nodePositions[sourceId];
              return sum + (pos ? pos.x : 0);
            }, 0) / incomingEdges[a].length
          : startX;
      const bAvgSourceX =
        incomingEdges[b].length > 0
          ? incomingEdges[b].reduce((sum, sourceId) => {
              const pos = nodePositions[sourceId];
              return sum + (pos ? pos.x : 0);
            }, 0) / incomingEdges[b].length
          : startX;
      return aAvgSourceX - bAvgSourceX;
    });

    sortedNodes.forEach((nodeId, index) => {
      nodePositions[nodeId] = {
        x: startXForLevel + index * nodeSpacing,
        y: y,
      };
    });
  });

  return nodePositions;
};

const RELATIONSHIP_TYPES = [
  { type: 'implements', color: '#4CAF50', label: 'Implements' },
  { type: 'depends_on', color: '#FF9800', label: 'Depends On' },
  { type: 'uses', color: '#2196F3', label: 'Uses' },
  { type: 'related_to', color: '#9C27B0', label: 'Related To' },
  { type: 'part_of', color: '#F44336', label: 'Part Of' },
  { type: 'documents', color: '#00BCD4', label: 'Documents' },
];

const getEdgeColor = (type: string): string => {
  const relationship = RELATIONSHIP_TYPES.find((r) => r.type === type);
  return relationship?.color || '#888888';
};

const convertGraphToFlow = (graph: Graph): { nodes: Node[]; edges: Edge[] } => {
  const nodePositions = calculateLayout(graph.nodes, graph.edges);

  const nodes: Node[] = graph.nodes.map((node: ApiGraphNode) => {
    const position = nodePositions[node.id] || { x: 500, y: 400 };
    const { width, height } = getNodeSize(node.name);
    const color = stringToColor(node.category);

    return {
      id: node.id.toString(),
      type: 'default',
      position,
      data: {
        label: (
          <div className="mindmap-node-content">
            <div className="mindmap-node-title">{node.name}</div>
            <div className="mindmap-node-category">{node.category}</div>
          </div>
        ),
      },
      style: {
        background: color,
        color: '#ffffff',
        border: '3px solid rgba(255, 255, 255, 0.8)',
        borderRadius: '12px',
        padding: '12px 16px',
        fontSize: '15px',
        fontWeight: '600',
        width: `${width}px`,
        height: `${height}px`,
        textAlign: 'center',
        boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15), 0 2px 4px rgba(0, 0, 0, 0.1)',
        cursor: 'pointer',
        transition: 'all 0.2s ease',
      },
    };
  });

  const edges: Edge[] = graph.edges.map((edge) => {
    const edgeColor = getEdgeColor(edge.type);

    return {
      id: edge.id.toString(),
      source: edge.source.toString(),
      target: edge.target.toString(),
      type: 'straight',
      animated: false,
      label: edge.type.replace('_', ' '),
      data: {
        relationshipType: edge.type,
        relationshipId: edge.id,
      },
      labelStyle: {
        fontSize: '9px',
        fontWeight: '600',
        fill: '#333333',
      },
      labelBgStyle: {
        fill: 'rgba(255, 255, 255, 0.95)',
        fillOpacity: 0.95,
        stroke: edgeColor,
        strokeWidth: 1,
        rx: 4,
        ry: 4,
      },
      style: {
        stroke: edgeColor,
        strokeWidth: 3,
        opacity: 0.7,
      },
      markerEnd: {
        type: MarkerType.ArrowClosed,
        color: edgeColor,
        width: 20,
        height: 20,
      },
    };
  });

  return { nodes, edges };
};

const MindMap: React.FC<MindMapProps> = ({ graph, onNodeClick, onGraphUpdate }) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [showLegend, setShowLegend] = useState(true);
  const [pendingConnection, setPendingConnection] = useState<Connection | null>(null);
  const [showTypeDialog, setShowTypeDialog] = useState(false);
  const [selectedEdge, setSelectedEdge] = useState<Edge | null>(null);
  const [showEdgeDialog, setShowEdgeDialog] = useState(false);
  const [editingRelationshipType, setEditingRelationshipType] = useState<string>('');

  useEffect(() => {
    const { nodes: flowNodes, edges: flowEdges } = convertGraphToFlow(graph);
    setNodes(flowNodes);
    setEdges(flowEdges);
  }, [graph, setNodes, setEdges]);

  const handleCreateRelationship = useCallback(
    async (relationshipType: string) => {
      if (!pendingConnection) return;

      const sourceId = parseInt(pendingConnection.source!);
      const targetId = parseInt(pendingConnection.target!);

      if (isNaN(sourceId) || isNaN(targetId)) {
        console.error('Invalid node IDs for connection');
        setPendingConnection(null);
        setShowTypeDialog(false);
        return;
      }

      const edgeType = 'straight';

      const tempEdge: Edge = {
        id: `temp-${Date.now()}`,
        source: pendingConnection.source!,
        target: pendingConnection.target!,
        sourceHandle: pendingConnection.sourceHandle || null,
        targetHandle: pendingConnection.targetHandle || null,
        type: edgeType,
        animated: false,
        style: {
          stroke: '#888888',
          strokeWidth: 3,
          opacity: 0.7,
        },
        markerEnd: {
          type: MarkerType.ArrowClosed,
          color: '#888888',
          width: 20,
          height: 20,
        },
      };

      setEdges((eds) => [...eds, tempEdge]);
      setShowTypeDialog(false);

      try {
        const relationship = await relationshipsApi.create({
          source_term_id: sourceId,
          target_term_id: targetId,
          type: relationshipType,
          description: '',
        });

        const edgeColor = getEdgeColor(relationship.type);
        const newEdge: Edge = {
          id: relationship.id.toString(),
          source: pendingConnection.source!,
          target: pendingConnection.target!,
          sourceHandle: pendingConnection.sourceHandle || null,
          targetHandle: pendingConnection.targetHandle || null,
          type: edgeType,
          animated: false,
          label: relationship.type.replace('_', ' '),
          data: {
            relationshipType: relationship.type,
            relationshipId: relationship.id,
          },
          labelStyle: {
            fontSize: '9px',
            fontWeight: '600',
            fill: '#333333',
          },
          labelBgStyle: {
            fill: 'rgba(255, 255, 255, 0.95)',
            fillOpacity: 0.95,
            stroke: edgeColor,
            strokeWidth: 1,
            rx: 4,
            ry: 4,
          },
          style: {
            stroke: edgeColor,
            strokeWidth: 3,
            opacity: 0.7,
          },
          markerEnd: {
            type: MarkerType.ArrowClosed,
            color: edgeColor,
            width: 20,
            height: 20,
          },
        };

        setEdges((eds) => {
          const filtered = eds.filter((e) => e.id !== tempEdge.id);
          return [...filtered, newEdge];
        });
      } catch (error) {
        console.error('Failed to create relationship:', error);
        setEdges((eds) => eds.filter((e) => e.id !== tempEdge.id));
        alert('Failed to create relationship. Please try again.');
      } finally {
        setPendingConnection(null);
      }
    },
    [setEdges, pendingConnection]
  );

  const onConnect = useCallback(
    (params: Connection) => {
      if (!params.source || !params.target) {
        return;
      }

      setPendingConnection(params);
      setShowTypeDialog(true);
    },
    []
  );

  const handleNodeClick = useCallback(
    (_event: React.MouseEvent, node: Node) => {
      onNodeClick(parseInt(node.id));
    },
    [onNodeClick]
  );

  const handleEdgeClick = useCallback(
    (_event: React.MouseEvent, edge: Edge) => {
      setSelectedEdge(edge);
      setEditingRelationshipType(edge.data?.relationshipType || '');
      setShowEdgeDialog(true);
    },
    []
  );

  const handleDeleteRelationship = useCallback(async () => {
    if (!selectedEdge?.data?.relationshipId) return;

    try {
      await relationshipsApi.delete(selectedEdge.data.relationshipId);
      setShowEdgeDialog(false);
      setSelectedEdge(null);
      if (onGraphUpdate) {
        onGraphUpdate();
      } else {
        window.location.reload();
      }
    } catch (error) {
      console.error('Failed to delete relationship:', error);
      alert('Не удалось удалить связь. Попробуйте снова.');
    }
  }, [selectedEdge, onGraphUpdate]);

  const handleEditRelationship = useCallback(async () => {
    if (!selectedEdge?.data?.relationshipId || !editingRelationshipType) return;

    const sourceId = parseInt(selectedEdge.source);
    const targetId = parseInt(selectedEdge.target);

    if (isNaN(sourceId) || isNaN(targetId)) {
      alert('Ошибка: неверные ID узлов');
      return;
    }

    try {
      await relationshipsApi.delete(selectedEdge.data.relationshipId);
      await relationshipsApi.create({
        source_term_id: sourceId,
        target_term_id: targetId,
        type: editingRelationshipType,
        description: '',
      });
      setShowEdgeDialog(false);
      setSelectedEdge(null);
      if (onGraphUpdate) {
        onGraphUpdate();
      } else {
        window.location.reload();
      }
    } catch (error) {
      console.error('Failed to update relationship:', error);
      alert('Не удалось обновить связь. Попробуйте снова.');
    }
  }, [selectedEdge, editingRelationshipType, onGraphUpdate]);

  const nodeTypes = useMemo(() => ({}), []);

  return (
    <div className="mindmap-container">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={handleNodeClick}
        onEdgeClick={handleEdgeClick}
        nodeTypes={nodeTypes}
        nodesDraggable={true}
        elementsSelectable={false}
        fitView
        fitViewOptions={{ padding: 0.2, maxZoom: 1.5 }}
        minZoom={0.3}
        maxZoom={2}
        defaultEdgeOptions={{
          type: 'straight',
          animated: false,
          style: {
            strokeWidth: 3,
            stroke: '#888888',
            opacity: 0.7,
          },
          markerEnd: {
            type: MarkerType.ArrowClosed,
            color: '#888888',
            width: 20,
            height: 20,
          },
        }}
        attributionPosition="bottom-left"
      >
        <Controls
          showInteractive={false}
          style={{
            backgroundColor: 'rgba(255, 255, 255, 0.9)',
            border: '1px solid #e0e0e0',
            borderRadius: '8px',
          }}
        />
        <Background
          variant={BackgroundVariant.Dots}
          gap={20}
          size={1.5}
          color="#e0e0e0"
        />
      </ReactFlow>

      {showLegend && (
        <div className="mindmap-legend">
          <div className="mindmap-legend-header">
            <h3>Relationship Types</h3>
            <button
              className="mindmap-legend-toggle"
              onClick={() => setShowLegend(false)}
              aria-label="Hide legend"
            >
              ×
            </button>
          </div>
          <div className="mindmap-legend-items">
            {RELATIONSHIP_TYPES.map((rel) => (
              <div key={rel.type} className="mindmap-legend-item">
                <div
                  className="mindmap-legend-color"
                  style={{ backgroundColor: rel.color }}
                />
                <span className="mindmap-legend-label">{rel.label}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {!showLegend && (
        <button
          className="mindmap-legend-show-button"
          onClick={() => setShowLegend(true)}
          title="Show legend"
        >
          Legend
        </button>
      )}

      {showTypeDialog && (
        <div className="mindmap-dialog-overlay" onClick={() => {
          setShowTypeDialog(false);
          setPendingConnection(null);
        }}>
          <div className="mindmap-dialog" onClick={(e) => e.stopPropagation()}>
            <h3>Select Relationship Type</h3>
            <div className="mindmap-dialog-options">
              {RELATIONSHIP_TYPES.map((rel) => (
                <button
                  key={rel.type}
                  className="mindmap-dialog-option"
                  onClick={() => handleCreateRelationship(rel.type)}
                  style={{ borderLeftColor: rel.color }}
                >
                  <div
                    className="mindmap-dialog-option-color"
                    style={{ backgroundColor: rel.color }}
                  />
                  <span>{rel.label}</span>
                </button>
              ))}
            </div>
            <button
              className="mindmap-dialog-cancel"
              onClick={() => {
                setShowTypeDialog(false);
                setPendingConnection(null);
              }}
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {showEdgeDialog && selectedEdge && (
        <div className="mindmap-dialog-overlay" onClick={() => {
          setShowEdgeDialog(false);
          setSelectedEdge(null);
        }}>
          <div className="mindmap-dialog" onClick={(e) => e.stopPropagation()}>
            <h3>Edit Relationship</h3>
            <div style={{ marginBottom: '20px' }}>
              <label style={{ display: 'block', marginBottom: '8px', fontWeight: '500' }}>
                Relationship Type:
              </label>
              <div className="mindmap-dialog-options">
                {RELATIONSHIP_TYPES.map((rel) => (
                  <button
                    key={rel.type}
                    className={`mindmap-dialog-option ${editingRelationshipType === rel.type ? 'selected' : ''}`}
                    onClick={() => setEditingRelationshipType(rel.type)}
                    style={{ borderLeftColor: rel.color }}
                  >
                    <div
                      className="mindmap-dialog-option-color"
                      style={{ backgroundColor: rel.color }}
                    />
                    <span>{rel.label}</span>
                  </button>
                ))}
              </div>
            </div>
            <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end' }}>
              <button
                className="mindmap-dialog-cancel"
                onClick={() => {
                  setShowEdgeDialog(false);
                  setSelectedEdge(null);
                }}
              >
                Cancel
              </button>
              <button
                className="mindmap-dialog-delete"
                onClick={handleDeleteRelationship}
                style={{
                  backgroundColor: '#f44336',
                  color: 'white',
                  border: 'none',
                  padding: '8px 16px',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  fontWeight: '500',
                }}
              >
                Delete
              </button>
              <button
                className="mindmap-dialog-save"
                onClick={handleEditRelationship}
                disabled={!editingRelationshipType}
                style={{
                  backgroundColor: '#2196F3',
                  color: 'white',
                  border: 'none',
                  padding: '8px 16px',
                  borderRadius: '4px',
                  cursor: editingRelationshipType ? 'pointer' : 'not-allowed',
                  fontWeight: '500',
                  opacity: editingRelationshipType ? 1 : 0.5,
                }}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MindMap;


