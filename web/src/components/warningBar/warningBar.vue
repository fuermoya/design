<script setup>
import { Close, Warning } from '@element-plus/icons-vue'
import { computed } from 'vue'

defineOptions({
  name: 'WarningBar',
})

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  type: {
    type: String,
    default: 'warning',
    validator: value => ['warning', 'info', 'success', 'error'].includes(value),
  },
  closable: {
    type: Boolean,
    default: false,
  },
  showIcon: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['close'])

const iconMap = {
  warning: Warning,
  info: 'InfoFilled',
  success: 'SuccessFilled',
  error: 'CircleCloseFilled',
}

const currentIcon = computed(() => iconMap[props.type])

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="warning-bar" :class="[`warning-bar--${type}`]">
    <div class="warning-bar__content">
      <el-icon v-if="showIcon" class="warning-bar__icon">
        <component :is="currentIcon" />
      </el-icon>
      <span class="warning-bar__text">{{ title }}</span>
    </div>
    <el-button
      v-if="closable"
      type="text"
      size="small"
      class="warning-bar__close"
      @click="handleClose"
    >
      <el-icon>
        <Close />
      </el-icon>
    </el-button>
  </div>
</template>

<style lang="scss" scoped>
// 暗色主题支持
html.dark {
  .warning-bar {
    &--warning {
      background-color: #451a03;
      border-color: #78350f;
      color: #fbbf24;
    }

    &--info {
      background-color: #1e3a8a;
      border-color: #1e40af;
      color: #93c5fd;
    }

    &--success {
      background-color: #14532d;
      border-color: #166534;
      color: #86efac;
    }

    &--error {
      background-color: #7f1d1d;
      border-color: #991b1b;
      color: #fca5a5;
    }
  }
}
</style>
