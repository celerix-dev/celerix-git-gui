<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as App from '../../../wailsjs/go/backend/App';
import type { SshKeyInfo } from '@/types/git.types';

defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const keyInfo = ref<SshKeyInfo | null>(null);
const loading = ref(false);
const generating = ref(false);
const error = ref<string | null>(null);
const copySuccess = ref(false);
const clearingCache = ref(false);

const loadKeyInfo = async () => {
  loading.value = true;
  error.value = null;
  try {
    keyInfo.value = await App.GetSshKeyInfo();
  } catch (err: any) {
    error.value = err.toString();
  } finally {
    loading.value = false;
  }
};

const generateKey = async () => {
  generating.value = true;
  error.value = null;
  try {
    keyInfo.value = await App.GenerateSshKey();
  } catch (err: any) {
    error.value = err.toString();
  } finally {
    generating.value = false;
  }
};

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text);
  copySuccess.value = true;
  setTimeout(() => {
    copySuccess.value = false;
  }, 2000);
};

onMounted(() => {
  loadKeyInfo();
});
</script>

<template>
  <div v-if="show" class="modal-backdrop fade show"></div>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-lg modal-dialog-centered" role="document">
      <div class="modal-content shadow-lg border-0">
        <div class="modal-header border-bottom-0">
          <h5 class="modal-title d-flex align-items-center">
            <i class="ti ti-settings me-2 text-primary"></i>
            Git Settings
          </h5>
          <button type="button" class="btn-close" @click="emit('close')" aria-label="Close"></button>
        </div>
        <div class="modal-body p-4">
          <div class="card shadow-sm mb-4 border-0 bg-light-subtle">
            <div class="card-header bg-transparent border-0 pt-3">
              <h6 class="mb-0 fw-bold">SSH Key Management</h6>
            </div>
            <div class="card-body">
              <div v-if="loading" class="text-center py-4">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Loading...</span>
                </div>
              </div>

              <div v-else-if="error" class="alert alert-danger">
                {{ error }}
              </div>

              <div v-else>
                <div v-if="!keyInfo?.has_key" class="text-center py-4">
                  <i class="ti ti-key fs-1 text-muted mb-3 d-block"></i>
                  <p>No SSH key (id_ed25519) found in your home directory.</p>
                  <button 
                    class="btn btn-primary" 
                    @click="generateKey" 
                    :disabled="generating"
                  >
                    <span v-if="generating" class="spinner-border spinner-border-sm me-2"></span>
                    Generate New SSH Key
                  </button>
                  <p class="small text-muted mt-3">
                    This will generate an Ed25519 key pair without a passphrase.
                  </p>
                </div>

                <div v-else>
                  <div class="mb-3">
                    <label class="form-label fw-bold small">SSH Public Key</label>
                    <div class="position-relative">
                      <textarea 
                        class="form-control font-monospace smaller bg-body" 
                        rows="5" 
                        readonly 
                        :value="keyInfo.public_key"
                        style="padding-right: 80px;"
                      ></textarea>
                      <button 
                        class="btn btn-sm btn-primary position-absolute top-0 end-0 m-2"
                        @click="copyToClipboard(keyInfo.public_key)"
                      >
                        <i :class="['ti', copySuccess ? 'ti-check' : 'ti-copy']"></i>
                        {{ copySuccess ? 'Copied!' : 'Copy' }}
                      </button>
                    </div>
                    <div class="form-text mt-2 smaller">
                      Copy this key and add it to your Git provider's settings (e.g., GitHub SSH and GPG keys).
                    </div>
                  </div>

                  <div class="mb-0">
                    <label class="form-label fw-bold small">Key Path</label>
                    <p class="text-muted smaller mb-0"><code>{{ keyInfo.path }}</code></p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="card shadow-sm border-0 bg-light-subtle">
            <div class="card-header bg-transparent border-0 pt-3">
              <h6 class="mb-0 fw-bold">Help & Context</h6>
            </div>
            <div class="card-body">
              <p class="small">To interact with remote repositories (GitHub, GitLab, etc.), you need to authenticate. SSH is the recommended method.</p>
              <ol class="small mb-0">
                <li>Generate an SSH key if you don't have one.</li>
                <li>Copy the public key shown above.</li>
                <li>Go to your Git provider (e.g., <a href="https://github.com/settings/keys" target="_blank">GitHub Settings</a>).</li>
                <li>Add the new SSH key.</li>
                <li>Now you can perform operations like Fetch, Pull, and Push.</li>
              </ol>
            </div>
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0 pb-4 pe-4">
          <button type="button" class="btn btn-secondary px-4" @click="emit('close')">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.smaller {
  font-size: 0.75rem;
}
</style>
