<script setup lang="ts">
import { computed } from 'vue';

interface CommitNode {
  hash: string;
  parentHashes: string[];
}

const props = defineProps<{
  commits: CommitNode[];
  rowHeight: number;
}>();

/**
 * Git Graph logic:
 * We need to assign each commit to a "lane".
 * When a commit has multiple parents, it's a merge.
 * When multiple commits have the same parent, it's a branch.
 */
const graphData = computed(() => {
  const nodes: any[] = [];
  const paths: any[] = [];
  
  let activeLanes: (string | null)[] = [];

  props.commits.forEach((commit, rowIndex) => {
    let laneIndex = activeLanes.indexOf(commit.hash);
    
    if (laneIndex === -1) {
      laneIndex = activeLanes.indexOf(null);
      if (laneIndex === -1) {
        laneIndex = activeLanes.length;
        activeLanes.push(commit.hash);
      } else {
        activeLanes[laneIndex] = commit.hash;
      }
    }

    const x = laneIndex * 15 + 10;
    const y = rowIndex * props.rowHeight + props.rowHeight / 2;
    
    nodes.push({
      x, y,
      color: getLaneColor(laneIndex),
      hash: commit.hash
    });

    const parents = commit.parentHashes || [];
    
    // Draw lines for lanes that are just passing through
    activeLanes.forEach((laneHash, lIdx) => {
      if (laneHash && laneHash !== commit.hash) {
        const lx = lIdx * 15 + 10;
        // Draw line from top of row to bottom of row
        paths.push({
          d: `M ${lx} ${rowIndex * props.rowHeight} L ${lx} ${(rowIndex + 1) * props.rowHeight}`,
          color: getLaneColor(lIdx)
        });
      } else if (laneHash === commit.hash) {
        const lx = lIdx * 15 + 10;
        // Draw line from top of row to the node
        paths.push({
          d: `M ${lx} ${rowIndex * props.rowHeight} L ${lx} ${y}`,
          color: getLaneColor(lIdx)
        });
      }
    });

    if (parents.length > 0) {
      // First parent connection
      const firstParentHash = parents[0];
      let firstParentLane = activeLanes.indexOf(firstParentHash);
      
      if (firstParentLane === -1) {
        // Parent not in any lane yet, reuse current lane or find empty lane
        let targetLane = laneIndex;
        activeLanes[targetLane] = firstParentHash;
        firstParentLane = targetLane;
      }
      
      const fx = firstParentLane * 15 + 10;

      // Every path MUST start at exactly the node's center or from the previous row's end.
      // And MUST end at the row boundary.
      
      if (firstParentLane === laneIndex) {
        // Straight line to next row
        paths.push({
          d: `M ${x} ${y} L ${x} ${(rowIndex + 1) * props.rowHeight}`,
          color: getLaneColor(laneIndex)
        });
      } else {
        // Curve to the parent lane. 
        // Starts at node (x, y), ends at (fx, (rowIndex + 1) * rowHeight)
        const startY = y;
        const endY = (rowIndex + 1) * props.rowHeight;
        const midY = (startY + endY) / 2;
        
        paths.push({
          d: `M ${x} ${startY} C ${x} ${midY}, ${fx} ${midY}, ${fx} ${endY}`,
          color: getLaneColor(laneIndex)
        });
      }

      // Handle additional parents (merges)
      for (let i = 1; i < parents.length; i++) {
        const parentHash = parents[i];
        let parentLane = activeLanes.indexOf(parentHash);
        if (parentLane === -1) {
          parentLane = activeLanes.indexOf(null);
          if (parentLane === -1) {
            parentLane = activeLanes.length;
            activeLanes.push(parentHash);
          } else {
            activeLanes[parentLane] = parentHash;
          }
        }
        
        const px = parentLane * 15 + 10;
        const startY = y;
        const endY = (rowIndex + 1) * props.rowHeight;
        const midY = (startY + endY) / 2;
        
        paths.push({
          d: `M ${x} ${startY} C ${x} ${midY}, ${px} ${midY}, ${px} ${endY}`,
          color: getLaneColor(parentLane)
        });
      }
      
      if (firstParentLane !== laneIndex) {
          activeLanes[laneIndex] = null;
      }
    } else {
      activeLanes[laneIndex] = null;
    }

    // Clean up trailing nulls to keep lane count reasonable
    while (activeLanes.length > 0 && activeLanes[activeLanes.length - 1] === null) {
      activeLanes.pop();
    }
  });

  return { nodes, paths, laneCount: Math.max(...nodes.map(n => n.x / 15), 0) + 1 };
});

function getLaneColor(index: number) {
  const colors = [
    '#3498db', '#e74c3c', '#2ecc71', '#f1c40f', '#9b59b6', 
    '#1abc9c', '#e67e22', '#34495e', '#d35400', '#c0392b'
  ];
  return colors[index % colors.length];
}

const svgWidth = computed(() => (graphData.value.laneCount + 1) * 15);
const svgHeight = computed(() => props.commits.length * props.rowHeight + 10);
</script>

<template>
  <svg :width="svgWidth" :height="svgHeight" class="git-graph-svg">
    <!-- Paths (Connectors) -->
    <path 
      v-for="(path, i) in graphData.paths" 
      :key="'p'+i"
      :d="path.d"
      :stroke="path.color"
      fill="none"
      stroke-width="2"
      stroke-linecap="round"
    />
    <!-- Nodes -->
    <circle 
      v-for="node in graphData.nodes" 
      :key="node.hash"
      :cx="node.x" :cy="node.y" 
      r="4"
      :fill="node.color"
      stroke="white"
      stroke-width="1"
    />
  </svg>
</template>

<style scoped>
.git-graph-svg {
  display: block;
}
</style>
