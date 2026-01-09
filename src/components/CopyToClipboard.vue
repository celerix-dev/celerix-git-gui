<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  text: string;
  color?: 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info' | 'light' | 'dark' | 'link';
  size?: 'sm' | 'lg' | 'md' | 'normal';
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  size: 'sm',
});

const isCopied = ref(false);

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(props.text);
    isCopied.value = true;
    setTimeout(() => {
      isCopied.value = false;
    }, 3000); // Reset the "Copied!" message after 3 seconds
  } catch (error) {
    console.error('Failed to copy: ', error);
  }
};
</script>

<template>
    <div :class="['btn', `btn-${color}`, (size !== 'md' && size !== 'normal') ? `btn-${size}` : '' , '']" @click="copyToClipboard">
        <i class="ti ti-copy" v-if="!isCopied"></i>
        <i class="ti ti-check" v-if="isCopied"></i>
    </div>
</template>