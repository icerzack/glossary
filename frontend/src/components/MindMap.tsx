import React, { useCallback, useEffect } from 'react';
import ReactFlow, {
  Node,
  Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  Connection,
  MarkerType,
  BackgroundVariant,
} from 'reactflow';
import 'reactflow/dist/style.css';
import type { Graph, GraphNode as ApiGraphNode } from '../types';
import { stringToColor } from '../utils/colorUtils';

interface MindMapProps {
  graph: Graph;
  onNodeClick: (nodeId: number) => void;
}

const getNodeColor = (category: string): string => {
  return stringToColor(category);
};

const convertGraphToFlow = (graph: Graph): { nodes: Node[]; edges: Edge[] } => {
  const nodes: Node[] = graph.nodes.map((node: ApiGraphNode, index: number) => {
    const angle = (2 * Math.PI * index) / graph.nodes.length;
    const radius = 300;
    const x = 500 + radius * Math.cos(angle);
    const y = 400 + radius * Math.sin(angle);

    return {
      id: node.id.toString(),
      type: 'default',
      position: { x, y },
      data: {
        label: node.name,
      },
      style: {
        background: getNodeColor(node.category),
        color: 'white',
        border: '2px solid #fff',
        borderRadius: '8px',
        padding: '10px',
        fontSize: '14px',
        fontWeight: 'bold',
        minWidth: '120px',
        textAlign: 'center',
      },
    };
  });

  const edges: Edge[] = graph.edges.map((edge) => ({
    id: edge.id.toString(),
    source: edge.source.toString(),
    target: edge.target.toString(),
    type: 'smoothstep',
    animated: true,
    label: edge.type,
    style: { stroke: '#888', strokeWidth: 2 },
    labelStyle: { fontSize: 10, fill: '#666' },
    markerEnd: {
      type: MarkerType.ArrowClosed,
      color: '#888',
    },
  }));

  return { nodes, edges };
};

const MindMap: React.FC<MindMapProps> = ({ graph, onNodeClick }) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  useEffect(() => {
    const { nodes: flowNodes, edges: flowEdges } = convertGraphToFlow(graph);
    setNodes(flowNodes);
    setEdges(flowEdges);
  }, [graph, setNodes, setEdges]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  const handleNodeClick = useCallback(
    (_event: React.MouseEvent, node: Node) => {
      onNodeClick(parseInt(node.id));
    },
    [onNodeClick]
  );

  return (
    <div style={{ width: '100%', height: '100%' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={handleNodeClick}
        fitView
        attributionPosition="bottom-left"
      >
        <Controls />
        <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
      </ReactFlow>
    </div>
  );
};

export default MindMap;

