<script setup lang="ts">
import FroalaEditor from 'froala-editor'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import 'froala-editor/css/froala_editor.pkgd.min.css'
import 'froala-editor/css/froala_style.min.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '请输入内容...',
  },
  height: {
    type: Number,
    default: 400,
  },
})

const emit = defineEmits(['update:modelValue', 'change'])

const editorRef = ref<HTMLElement>()
let editor: any = null
const isEditorReady = ref(false)

// 初始化编辑器
onMounted(() => {
  if (editorRef.value) {
    editor = new FroalaEditor(editorRef.value, {
      placeholderText: props.placeholder,
      height: props.height,
      toolbarButtons: [
        'bold',
        'italic',
        'underline',
        'strikeThrough',
        '|',
        'fontSize',
        'color',
        'backgroundColor',
        '|',
        'paragraphFormat',
        'align',
        'formatOL',
        'formatUL',
        '|',
        'outdent',
        'indent',
        '|',
        'insertLink',
        'insertImage',
        'insertVideo',
        'insertTable',
        '|',
        'undo',
        'redo',
        '|',
        'fullscreen',
      ],
      events: {
        contentChanged: () => {
          const content = editor?.html.get() || ''
          emit('update:modelValue', content)
          emit('change', content)
        },
        initialized: () => {
          // 设置初始值
          isEditorReady.value = true
          if (props.modelValue) {
            editor?.html.set(props.modelValue)
          }
        },
      },
      // 图片上传配置
      imageUploadURL: '/api/upload/image',
      imageUploadParams: {
        type: 'image',
      },
      // 视频上传配置
      videoUploadURL: '/api/upload/video',
      videoUploadParams: {
        type: 'video',
      },
      // 文件上传配置
      fileUploadURL: '/api/upload/file',
      fileUploadParams: {
        type: 'file',
      },
      // 字体大小选项
      fontSizeSelection: true,
      fontSize: ['8', '10', '12', '14', '16', '18', '20', '24', '30', '36', '48', '60', '72', '96'],
      // 段落格式选项
      paragraphFormatSelection: true,
      paragraphFormat: {
        N: 'Normal',
        h1: 'Heading 1',
        h2: 'Heading 2',
        h3: 'Heading 3',
        h4: 'Heading 4',
        h5: 'Heading 5',
        h6: 'Heading 6',
        pre: 'Preformatted',
      },
      // 颜色选择器
      colorsBackground: [
        '#FFFFFF',
        '#000000',
        '#eeece1',
        '#1f497d',
        '#4f81bd',
        '#c05050',
        '#9bbb59',
        '#8064a2',
        '#4bacc6',
        '#f79646',
        '#f2f2f2',
        '#7f7f7f',
        '#ddd9c3',
        '#c6d9f0',
        '#dbe5f1',
        '#f2dcdb',
        '#ebf1dd',
        '#e5e0ec',
        '#dbeef3',
        '#fde9d9',
        '#d8d8d8',
        '#595959',
        '#c4bd97',
        '#8db3e2',
        '#b8cce4',
        '#e6b8b7',
        '#d7e3bc',
        '#ccc1d9',
        '#b7dde8',
        '#fbd5b5',
        '#bfbfbf',
        '#404040',
        '#938953',
        '#548dd4',
        '#95b3d7',
        '#d99694',
        '#c3d69b',
        '#b2a2c7',
        '#92cddc',
        '#fac08f',
        '#a5a5a5',
        '#262626',
        '#494429',
        '#17365d',
        '#366092',
        '#953734',
        '#76923c',
        '#5f497a',
        '#31869b',
        '#e36c09',
      ],
      colorsText: [
        '#FFFFFF',
        '#000000',
        '#eeece1',
        '#1f497d',
        '#4f81bd',
        '#c05050',
        '#9bbb59',
        '#8064a2',
        '#4bacc6',
        '#f79646',
        '#f2f2f2',
        '#7f7f7f',
        '#ddd9c3',
        '#c6d9f0',
        '#dbe5f1',
        '#f2dcdb',
        '#ebf1dd',
        '#e5e0ec',
        '#dbeef3',
        '#fde9d9',
        '#d8d8d8',
        '#595959',
        '#c4bd97',
        '#8db3e2',
        '#b8cce4',
        '#e6b8b7',
        '#d7e3bc',
        '#ccc1d9',
        '#b7dde8',
        '#fbd5b5',
        '#bfbfbf',
        '#404040',
        '#938953',
        '#548dd4',
        '#95b3d7',
        '#d99694',
        '#c3d69b',
        '#b2a2c7',
        '#92cddc',
        '#fac08f',
        '#a5a5a5',
        '#262626',
        '#494429',
        '#17365d',
        '#366092',
        '#953734',
        '#76923c',
        '#5f497a',
        '#31869b',
        '#e36c09',
      ],
      // 表格配置
      tableCellStyles: {
        red: 'Red',
        green: 'Green',
        blue: 'Blue',
      },
      // 链接配置
      linkAlwaysBlank: true,
      linkAutoPrefix: 'http://',
      // 其他配置
      attribution: false,
      charCounterCount: true,
      charCounterMax: 10000,
      quickInsertButtons: ['image', 'table', 'ol', 'ul', 'hr'],
    })
  }
})

// 监听外部值变化
watch(() => props.modelValue, (newVal: string) => {
  if (isEditorReady.value && editor && editor.html && newVal !== editor.html.get()) {
    editor.html.set(newVal)
  }
})

// 监听placeholder变化
watch(() => props.placeholder, (newVal: string) => {
  if (isEditorReady.value && editor && editor.placeholder) {
    editor.placeholder.set(newVal)
  }
})

// 监听height变化
watch(() => props.height, (newVal: number) => {
  if (isEditorReady.value && editor && editor.height) {
    editor.height.set(newVal)
  }
})

// 组件销毁时清理编辑器
onBeforeUnmount(() => {
  if (editor) {
    editor.destroy()
    editor = null
    isEditorReady.value = false
  }
})

// 暴露编辑器实例
defineExpose({
  getEditor: () => editor,
  getContent: () => (isEditorReady.value && editor?.html?.get()) || '',
  setContent: (content: string) => {
    if (isEditorReady.value && editor && editor.html) {
      editor.html.set(content)
    }
  },
  insertHTML: (html: string) => {
    if (isEditorReady.value && editor && editor.html) {
      editor.html.insert(html)
    }
  },
  insertImage: (url: string, alt?: string) => {
    if (isEditorReady.value && editor && editor.image) {
      editor.image.insert(url, null, null, null, null, alt)
    }
  },
  insertLink: (url: string, text?: string) => {
    if (isEditorReady.value && editor && editor.link) {
      editor.link.insert(url, text)
    }
  },
})
</script>

<template>
  <div class="html-editor-container">
    <div ref="editorRef" class="froala-editor" />
  </div>
</template>

<style scoped>
.html-editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  width: 100%;
}

.froala-editor {
  min-height: 400px;
}

/* 自定义Froala编辑器样式 */
:deep(.fr-box) {
  border: none;
}

:deep(.fr-toolbar) {
  border-bottom: 1px solid #dcdfe6;
  background-color: #f5f7fa;
}

:deep(.fr-toolbar .fr-btn) {
  color: #606266;
}

:deep(.fr-toolbar .fr-btn:hover) {
  background-color: #e6f7ff;
  color: #409eff;
}

:deep(.fr-toolbar .fr-btn.fr-active) {
  background-color: #409eff;
  color: white;
}

:deep(.fr-wrapper) {
  background-color: white;
}

:deep(.fr-element) {
  padding: 16px;
  min-height: 300px;
  line-height: 1.6;
}

:deep(.fr-element:focus) {
  outline: none;
}

/* 暗色主题支持 */
html.dark :deep(.fr-toolbar) {
  background-color: #2b2b2c;
  border-bottom-color: #4c4d4f;
}

html.dark :deep(.fr-toolbar .fr-btn) {
  color: #e5eaf3;
}

html.dark :deep(.fr-toolbar .fr-btn:hover) {
  background-color: #3a3b3c;
  color: #409eff;
}

html.dark :deep(.fr-wrapper) {
  background-color: #2b2b2c;
}

html.dark :deep(.fr-element) {
  background-color: #2b2b2c;
  color: #e5eaf3;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .froala-editor {
    min-height: 300px;
  }

  :deep(.fr-toolbar) {
    padding: 4px;
  }

  :deep(.fr-toolbar .fr-btn) {
    padding: 6px 8px;
    font-size: 12px;
  }
}
</style>
