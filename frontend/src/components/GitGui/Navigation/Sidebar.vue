<script setup lang="ts">

import { ref, computed } from 'vue';
import type { RepoTab } from "@/types/git.types";
import { backend } from "../../../../wailsjs/go/models";
import RepoStats = backend.RepoStats;

const props = defineProps<{
  recentRepos: {name: string, path: string}[];
  activeTab: RepoTab | null;
  currentRepoStats: RepoStats | null;
}>();

const emit = defineEmits<{
  (e: 'openSettings'): void;
  (e: 'openRepo'): void;
  (e: 'openRecentRepo', repo: {name: string, path: string}): void;
  (e: 'selectVerticalTab', tab: 'info' | 'local-changes' | 'commit' | 'placeholder1' | 'placeholder2'): void;
  (e: 'checkoutBranch', branchName: string, isRemote: boolean): void;
  (e: 'newBranch', fromBranch: string): void;
  (e: 'newTag', fromBranch: string): void;
}>();

interface BranchNode {
  name: string;
  fullName: string;
  isFolder: boolean;
  children: BranchNode[];
}

const collapsedFolders = ref<Record<string, boolean>>({});

const toggleFolder = (folderPath: string) => {
  collapsedFolders.value = {
    ...collapsedFolders.value,
    [folderPath]: !collapsedFolders.value[folderPath]
  };
};

const buildBranchTree = (branches: string[]): BranchNode[] => {
  const root: BranchNode[] = [];

  branches.forEach(branch => {
    const parts = branch.split('/');
    let currentLevel = root;
    let currentPath = '';

    parts.forEach((part, index) => {
      currentPath = currentPath ? `${currentPath}/${part}` : part;
      const isLast = index === parts.length - 1;
      
      let node = currentLevel.find(n => n.name === part);
      
      if (!node) {
        node = {
          name: part,
          fullName: isLast ? branch : currentPath,
          isFolder: !isLast,
          children: []
        };
        currentLevel.push(node);
      } else if (!isLast) {
        node.isFolder = true;
      }
      
      currentLevel = node.children;
    });
  });

  const sortNodes = (nodes: BranchNode[]): BranchNode[] => {
    return nodes.sort((a, b) => {
      // Branches without slashes (top-level leaf nodes in the tree context, 
      // but the requirement says "all branches that do not have a slash in them should be presented first")
      // In our tree, these are nodes where fullName doesn't contain a slash AND is not a folder.
      
      const aIsSimpleBranch = !a.isFolder && !a.fullName.includes('/');
      const bIsSimpleBranch = !b.isFolder && !b.fullName.includes('/');
      
      if (aIsSimpleBranch && !bIsSimpleBranch) return -1;
      if (!aIsSimpleBranch && bIsSimpleBranch) return 1;
      
      // Next, folders vs leaf nodes at this level
      if (!a.isFolder && b.isFolder) return -1;
      if (a.isFolder && !b.isFolder) return 1;
      
      return a.name.localeCompare(b.name);
    });
  };

  const finalize = (nodes: BranchNode[]) => {
    sortNodes(nodes);
    nodes.forEach(node => {
      if (node.children.length > 0) {
        finalize(node.children);
      }
    });
  };

  finalize(root);
  return root;
};

const localBranchTree = computed(() => {
  return buildBranchTree(props.currentRepoStats?.branches || []);
});

const remoteBranchTrees = computed(() => {
  const trees: Record<string, BranchNode[]> = {};
  props.currentRepoStats?.remotes?.forEach(remote => {
    const branches = (remote.branches || [])
      .filter(b => !b.endsWith('/HEAD'))
      .map(b => b.startsWith(remote.name + '/') ? b.substring(remote.name.length + 1) : b);
    trees[remote.name] = buildBranchTree(branches);
  });
  return trees;
});

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  branchName: '',
  isRemote: false
});

const showContextMenu = (event: MouseEvent, branchName: string, isRemote: boolean) => {
  event.preventDefault();
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    branchName,
    isRemote
  };
  
  const closeMenu = () => {
    contextMenu.value.show = false;
    document.removeEventListener('click', closeMenu);
  };
  document.addEventListener('click', closeMenu);
};

const handleCheckout = () => {
  if (contextMenu.value.isRemote) {
    // If it's on remote and it's already the active branch on local, it does nothing
    // Remote branch name can be origin/main or origin/feature/branch
    const parts = contextMenu.value.branchName.split('/');
    const localName = parts.slice(1).join('/');
    if (localName === props.currentRepoStats?.currentBranch) {
      contextMenu.value.show = false;
      return;
    }
  } else if (contextMenu.value.branchName === props.currentRepoStats?.currentBranch) {
    contextMenu.value.show = false;
    return;
  }
  
  emit('checkoutBranch', contextMenu.value.branchName, contextMenu.value.isRemote);
  contextMenu.value.show = false;
};

const collapsed = ref({
  branches: false, // Branches default opened (false means NOT collapsed)
  remotes: true,
  tags: true,
  stashes: true
});

const collapsedRemotes = ref<Record<string, boolean>>({
  'origin': false // origin default opened
});

const toggleSection = (section: keyof typeof collapsed.value) => {
  collapsed.value[section] = !collapsed.value[section];
};

const toggleRemote = (remoteName: string) => {
  const currentValue = collapsedRemotes.value[remoteName] !== undefined ? collapsedRemotes.value[remoteName] : false;
  collapsedRemotes.value = {
    ...collapsedRemotes.value,
    [remoteName]: !currentValue
  };
};

</script>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';

const BranchTreeNode = defineComponent({
  name: 'BranchTreeNode',
  props: {
    node: {
      type: Object as PropType<any>, // Changed to any to avoid type issue in template script
      required: true
    },
    currentBranch: String,
    isRemote: Boolean,
    remotePrefix: {
      type: String,
      default: ''
    },
    collapsedFolders: {
      type: Object as PropType<Record<string, boolean>>,
      required: true
    },
    level: {
      type: Number,
      default: 0
    }
  },
  emits: ['toggle-folder', 'contextmenu', 'checkout'],
  setup(props, { emit }) {
    const handleContext = (e: MouseEvent) => {
      if (!props.node.isFolder) {
        const fullBranchName = props.isRemote ? props.remotePrefix + props.node.fullName : props.node.fullName;
        emit('contextmenu', e, fullBranchName, props.isRemote);
      }
    };

    const handleCheckout = () => {
      if (!props.node.isFolder) {
        const fullBranchName = props.isRemote ? props.remotePrefix + props.node.fullName : props.node.fullName;
        emit('checkout', fullBranchName);
      }
    };

    const toggle = () => {
      if (props.node.isFolder) {
        emit('toggle-folder', props.node.fullName);
      }
    };

    return {
      handleContext,
      handleCheckout,
      toggle
    };
  },
  template: `
    <div class="branch-tree-node">
      <div 
        v-if="node.isFolder"
        class="list-group-item d-flex align-items-center border-0 cursor-pointer py-1"
        :style="{ paddingLeft: (level * 1) + 'rem' }"
        @click="toggle"
      >
        <i :class="['ti ms-3 me-2', collapsedFolders[node.fullName] ? 'ti-folder opacity-75' : 'ti-folder-open']"></i>
        <span class="text-truncate">{{ node.name }}</span>
      </div>
      
      <a v-else
         href="javascript:void(0);"
         :class="['list-group-item list-group-item-action d-flex align-items-center border-0 py-1', { 'fw-bold active-branch': (node.fullName === currentBranch) }]"
         :style="{ paddingLeft: (level * 1 + 1.2) + 'rem' }"
         @contextmenu="handleContext"
         @dblclick="handleCheckout"
      >
        <i :class="['ti me-2', (node.fullName === currentBranch) ? 'ti-check' : 'ti-git-branch']"></i>
        <span class="text-truncate">{{ node.name }}</span>
      </a>

      <div v-if="node.isFolder && !collapsedFolders[node.fullName]" class="branch-folder-children">
        <BranchTreeNode 
          v-for="child in node.children" 
          :key="child.fullName"
          :node="child"
          :current-branch="currentBranch"
          :is-remote="isRemote"
          :remote-prefix="remotePrefix"
          :collapsed-folders="collapsedFolders"
          :level="level + 1"
          @toggle-folder="(path) => $emit('toggle-folder', path)"
          @contextmenu="(e, name, remote) => $emit('contextmenu', e, name, remote)"
          @checkout="(name) => $emit('checkout', name)"
        />
      </div>
    </div>
  `
});

export default {
  components: {
    BranchTreeNode
  }
}
</script>

<template>
  <nav class="sidebar offcanvas-start offcanvas-md" tabindex="-1" id="donkers-sidebar">
    <div class="offcanvas-body p-0 d-flex flex-column">
      <!-- Top Part: Management Actions & Active Repo -->
      <div class="sidebar-management-section p-2 border-bottom">
        <div class="sidebar-group">
          <h6 class="sidebar-header">Actions</h6>
          <div class="list-group list-group-flush">
            <a href="javascript:void(0);" class="list-group-item list-group-item-action d-flex align-items-center border-0" @click.prevent="emit('openRepo')">
              <i class="ti ti-folder-open me-2"></i>
              <span>Open Repository</span>
            </a>
            <a href="javascript:void(0);" class="list-group-item list-group-item-action d-flex align-items-center border-0" @click.prevent="emit('openSettings')">
              <i class="ti ti-settings me-2"></i>
              <span>Settings</span>
            </a>
          </div>

          <template v-if="props.activeTab">
            <h6 class="sidebar-header text-truncate mt-2" :title="props.activeTab.name">{{ props.activeTab.name }}</h6>
            <div class="list-group list-group-flush">
              <a href="javascript:void(0);" :class="['list-group-item list-group-item-action d-flex align-items-center border-0', { active: props.activeTab.activeVerticalTab === 'local-changes' }]" @click="emit('selectVerticalTab', 'local-changes')">
                <i class="ti ti-edit me-2"></i>
                <span>Local changes {{ props.currentRepoStats?.modifiedFiles?.length ? `(${props.currentRepoStats.modifiedFiles.length})` : '' }}</span>
              </a>
              <a href="javascript:void(0);" :class="['list-group-item list-group-item-action d-flex align-items-center border-0', { active: props.activeTab.activeVerticalTab === 'commit' }]" @click="emit('selectVerticalTab', 'commit')">
                <i class="ti ti-history me-2"></i>
                <span>All Commits</span>
              </a>
            </div>
          </template>
        </div>
      </div>

      <!-- Bottom Part: Content Navigation -->
      <ul class="sidebar-nav flex-grow-1 overflow-auto py-2" data-key="dashboard" aria-expanded="true">
        <li v-if="props.activeTab">
          <div class="px-3 mb-3 mt-2">
            <input type="text" class="form-control form-control-sm filter-input" placeholder="Filter"/>
          </div>
          <h6 class="sidebar-header collapsible" @click="toggleSection('branches')">
            <i :class="['ti ti-chevron-right transition-icon me-2', { 'rotate-90': !collapsed.branches }]"></i>
            <span>Branches</span>
          </h6>
          <div v-if="!collapsed.branches" class="list-group list-group-flush" style="padding-left: 17px;">
            <template v-for="node in localBranchTree" :key="node.fullName">
              <BranchTreeNode 
                :node="node" 
                :current-branch="props.currentRepoStats?.currentBranch"
                :is-remote="false"
                :collapsed-folders="collapsedFolders"
                @toggle-folder="toggleFolder"
                @contextmenu="showContextMenu"
                @checkout="(name) => emit('checkoutBranch', name, false)"
              />
            </template>
            <div v-if="!props.currentRepoStats?.branches?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No branches found</span>
            </div>
          </div>
        </li>
        <li v-if="props.activeTab">
          <h6 class="sidebar-header collapsible" @click="toggleSection('remotes')">
            <i :class="['ti ti-chevron-right transition-icon me-2', { 'rotate-90': !collapsed.remotes }]"></i>
            <span>Remotes</span>
          </h6>
          <div v-if="!collapsed.remotes" class="list-group list-group-flush" style="padding-left: 17px;">
            <template v-for="remote in props.currentRepoStats?.remotes || []" :key="remote.name">
              <div 
                class="list-group-item d-flex align-items-center border-0 cursor-pointer py-1"
                @click="toggleRemote(remote.name)"
              >
                <i :class="['ti ti-chevron-right transition-icon me-2', { 'rotate-90': !collapsedRemotes[remote.name] }]"></i>
                <i class="ti ti-cloud me-2"></i>
                <span class="text-truncate">{{ remote.name }}</span>
              </div>
              <div v-if="!collapsedRemotes[remote.name]" class="remote-branches ms-3">
                <template v-for="node in remoteBranchTrees[remote.name]" :key="node.fullName">
                  <BranchTreeNode 
                    :node="node" 
                    :current-branch="props.currentRepoStats?.currentBranch"
                    :is-remote="true"
                    :remote-prefix="remote.name + '/'"
                    :collapsed-folders="collapsedFolders"
                    @toggle-folder="toggleFolder"
                    @contextmenu="showContextMenu"
                    @checkout="(name) => emit('checkoutBranch', name, true)"
                  />
                </template>
                <div v-if="!remote.branches?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
                  <span class="ms-4">No remote branches</span>
                </div>
              </div>
            </template>
            <div v-if="!props.currentRepoStats?.remotes?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No remotes found</span>
            </div>
          </div>
        </li>
        <li v-if="props.activeTab">
          <h6 class="sidebar-header collapsible" @click="toggleSection('tags')">
            <i :class="['ti ti-chevron-right transition-icon me-2', { 'rotate-90': !collapsed.tags }]"></i>
            <span>Tags</span>
          </h6>
          <div v-if="!collapsed.tags" class="list-group list-group-flush" style="padding-left: 17px;">
            <a v-for="tag in props.currentRepoStats?.tags || []"
               :key="tag"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center border-0"
            >
              <i class="ti ti-tag me-2"></i>
              <span class="text-truncate">{{ tag }}</span>
            </a>
            <div v-if="!props.currentRepoStats?.tags?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No tags found</span>
            </div>
          </div>
        </li>
        <li v-if="props.activeTab">
          <h6 class="sidebar-header collapsible" @click="toggleSection('stashes')">
            <i :class="['ti ti-chevron-right transition-icon me-2', { 'rotate-90': !collapsed.stashes }]"></i>
            <span>Stashes</span>
          </h6>
          <div v-if="!collapsed.stashes" class="list-group list-group-flush" style="padding-left: 17px;">
            <a v-for="(stash, index) in props.currentRepoStats?.stashes || []"
               :key="index"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center border-0"
            >
              <i class="ti ti-archive me-2"></i>
              <span class="text-truncate">{{ stash }}</span>
            </a>
            <div v-if="!props.currentRepoStats?.stashes?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No stashes found</span>
            </div>
          </div>
        </li>
        <li>
          <hr class="sidebar-divider mx-3 mt-1 mb-2" v-if="props.activeTab">
          <h6 class="sidebar-header">Recent Repositories</h6>
          <div class="list-group list-group-flush">
            <a v-for="repo in props.recentRepos"
               :key="repo.path"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center border-0"
               @click="emit('openRecentRepo', repo)"
               :title="repo.path"
            >
              <i class="ti ti-folder me-2"></i>
              <span class="text-truncate">{{ repo.name }}</span>
            </a>
          </div>
        </li>
      </ul>
    </div>
  </nav>

  <!-- Context Menu -->
  <div v-if="contextMenu.show" 
       class="dropdown-menu show position-fixed shadow-sm" 
       :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px', zIndex: 1050, minWidth: '200px' }"
  >
    <div class="dropdown-header border-bottom mb-1 py-1 text-truncate" style="max-width: 250px;">{{ contextMenu.branchName }}</div>
    <a class="dropdown-item d-flex align-items-center py-2" href="javascript:void(0);" @click="handleCheckout">
      <i class="ti ti-git-pull-request me-2"></i>
      <span>Check-out</span>
    </a>
    <div class="dropdown-divider my-1"></div>
    <a class="dropdown-item d-flex align-items-center py-2" href="javascript:void(0);" @click="emit('newBranch', contextMenu.branchName); contextMenu.show = false;">
      <i class="ti ti-git-branch me-2"></i>
      <span>New branch...</span>
    </a>
    <a class="dropdown-item d-flex align-items-center py-2" href="javascript:void(0);" @click="emit('newTag', contextMenu.branchName); contextMenu.show = false;">
      <i class="ti ti-tag me-2"></i>
      <span>New tag...</span>
    </a>
  </div>
</template>

<style scoped>
.sidebar-management-section {
  background-color: var(--bs-tertiary-bg);
}

.sidebar-header {
  height: 32px;
  display: flex;
  align-items: center;
  margin-bottom: 0;
  padding: 0 1rem;
}

.sidebar-divider {
  border-top: 1px solid var(--bs-border-color);
  opacity: 0.5;
  margin: 0.5rem 1rem;
}

.sidebar-header.collapsible {
  cursor: pointer;
  user-select: none;
}

.sidebar-header.collapsible:hover {
  background-color: var(--bs-tertiary-bg);
}

.transition-icon {
  width: 16px;
  height: 16px;
  display: inline-flex !important;
  align-items: center;
  justify-content: center;
  line-height: 1;
  transition: transform 0.2s ease;
  transform: rotate(0deg);
  transform-origin: center;
}

.rotate-90 {
  transform: rotate(90deg) !important;
}
.remote-branches .list-group-item {
  padding-left: 1rem !important; /* Adjusted for tree structure */
}
.branch-tree-node .list-group-item {
  height: 28px; /* Slightly smaller for tree items */
}

.list-group-item {
  border: none !important;
  height: 32px; /* Fixed height to prevent shifts */
  padding-left: 1.5rem !important; /* Increased indentation */
}

.list-group-item.active {
  border: none !important;
  margin-top: 0 !important;
}
.list-group-item i {
  width: 20px;
  text-align: center;
}
.filter-input {
  font-size: 0.8125rem;
  padding-top: 2px;
  padding-bottom: 2px;
  height: 28px;
}
</style>