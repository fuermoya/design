<script setup>
import { ElMessage } from 'element-plus'
import { ref, watch } from 'vue'
import request from '~/utils/request'

const props = defineProps({
  row: {
    default() {
      return {}
    },
    type: Object,
  },
})

const apiDefaultProps = ref({
  children: 'children',
  label: 'description',
})
const filterText = ref('')
const apiTreeData = ref([])
const apiTreeIds = ref([])
const activeUserId = ref('')
async function init() {
  const res2 = await request.post('/api/getAllApis')
  const apis = res2.data.apis
  apiTreeData.value = buildApiTree(apis)
  const res = await request.post('/casbin/getPolicyPathByAuthorityId', {
    authorityId: props.row.authorityId,
  })
  activeUserId.value = props.row.authorityId
  apiTreeIds.value = []
  res.data.paths && res.data.paths.forEach((item) => {
    apiTreeIds.value.push(`p:${item.path}m:${item.method}`)
  })
}

init()

const needConfirm = ref(false)
function nodeChange() {
  needConfirm.value = true
}
// 暴露给外层使用的切换拦截统一方法
function enterAndNext() {
  authApiEnter()
}

// 创建api树方法
function buildApiTree(apis) {
  const apiObj = {}
  apis
  && apis.forEach((item) => {
    item.onlyId = `p:${item.path}m:${item.method}`
    if (Object.prototype.hasOwnProperty.call(apiObj, item.apiGroup)) {
      apiObj[item.apiGroup].push(item)
    }
    else {
      Object.assign(apiObj, { [item.apiGroup]: [item] })
    }
  })
  const apiTree = []
  for (const key in apiObj) {
    const treeNode = {
      ID: key,
      description: `${key}组`,
      children: apiObj[key],
    }
    apiTree.push(treeNode)
  }
  return apiTree
}

// 关联关系确定
const apiTree = ref(null)
async function authApiEnter() {
  const checkArr = apiTree.value.getCheckedNodes(true)
  const casbinInfos = []
  checkArr && checkArr.forEach((item) => {
    const casbinInfo = {
      path: item.path,
      method: item.method,
    }
    casbinInfos.push(casbinInfo)
  })
  const res = await request.post('/casbin/updateCasbin', {
    authorityId: activeUserId.value,
    casbinInfos,
  })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: 'api设置成功' })
  }
}

defineExpose({
  needConfirm,
  enterAndNext,
})

function filterNode(value, data) {
  if (!value)
    return true
  return data.description.includes(value)
}
watch(filterText, (val) => {
  apiTree.value.filter(val)
})
</script>

<template>
  <div>
    <div class="sticky top-0.5 z-10 bg-white">
      <el-input
        v-model="filterText"
        class="w-3/5"
        placeholder="筛选"
      />
      <el-button
        class="float-right"
        type="primary"
        @click="authApiEnter"
      >
        确 定
      </el-button>
    </div>
    <div class="tree-content">
      <el-scrollbar>
        <el-tree
          ref="apiTree"
          :data="apiTreeData"
          :default-checked-keys="apiTreeIds"
          :props="apiDefaultProps"
          default-expand-all
          highlight-current
          node-key="onlyId"
          show-checkbox
          :filter-node-method="filterNode"
          @check="nodeChange"
        />
      </el-scrollbar>
    </div>
  </div>
</template>
