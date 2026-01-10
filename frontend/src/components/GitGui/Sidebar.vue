<script setup lang="ts">

const props = defineProps<{
  recentRepos: {name: string, path: string}[];
}>();

const emit = defineEmits<{
  (e: 'openSettings'): void;
  (e: 'openRepo'): void;
  (e: 'openRecentRepo', repo: {name: string, path: string}): void;
}>();

</script>

<template>
  <nav class="sidebar offcanvas-start offcanvas-md" tabindex="-1" id="donkers-sidebar">
    <div class="offcanvas-body">
      <div class="mb-3">
        <input type="text" class="form-control" placeholder="Search"/>
      </div>
      <ul class="sidebar-nav" data-key="dashboard" aria-expanded="true">
        <li>
          <h6 class="sidebar-header">Actions</h6>
        </li>
        <li class="nav-item list-group list-group-flush">
          <div class="list-group-item list-group-item-action ps-3 cursor-pointer" @click.prevent="emit('openRepo')"><i class="ti ti-folder-open me-1"></i> Open Repository</div>
        </li>
        <li class="nav-item list-group list-group-flush">
          <div class="list-group-item list-group-item-action ps-3 cursor-pointer" @click.prevent="emit('openSettings')"><i class="ti ti-settings me-1"></i> Settings</div>
        </li>
        <li>
          <h6 class="sidebar-header">Recent Repositories</h6>
          <div class="list-group list-group-flush">
            <a v-for="repo in props.recentRepos"
               :key="repo.path"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center py-1 border-0"
               @click="emit('openRecentRepo', repo)"
               :title="repo.path"
            >
              <i class="ti ti-folder me-2"></i>
              <span class="text-truncate small">{{ repo.name }}</span>
            </a>
          </div>
        </li>
      </ul>
    </div>
  </nav>
</template>

<style scoped>

</style>