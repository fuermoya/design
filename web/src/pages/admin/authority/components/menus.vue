<script setup>
import { ElMessage } from 'element-plus'
import { ref, watch } from 'vue'
import request from '~/utils/request'

defineOptions({
  name: 'Menus',
})

const props = defineProps({
  row: {
    default() {
      return {}
    },
    type: Object,
  },
})

const _emit = defineEmits(['changeRow'])
const filterText = ref('')
const menuTreeData = ref([])
const menuTreeIds = ref([])
const needConfirm = ref(false)
const menuDefaultProps = ref({
  children: 'children',
  label(data) {
    return data.meta.title
  },
})

async function init() {
  // 获取所有菜单树
  const res = await request.post('/menu/getBaseMenuTree')
  menuTreeData.value = res.data.menus
  const res1 = await request.post('/menu/getMenuAuthority', { authorityId: props.row.authorityId })
  const menus = res1.data.menus
  const arr = []

  // 递归检查菜单是否为叶子节点
  function isLeafNode(menuId, menuTree) {
    const menu = findMenuById(menuId, menuTree)
    if (!menu)
      return true
    return !menu.children || menu.children.length === 0
  }

  // 在菜单树中查找指定ID的菜单
  function findMenuById(menuId, menuTree) {
    for (const menu of menuTree) {
      if (menu.ID.toString() === menuId) {
        return menu
      }
      if (menu.children) {
        const found = findMenuById(menuId, menu.children)
        if (found)
          return found
      }
    }
    return null
  }

  // 只选中叶子节点
  menus.forEach((item) => {
    if (isLeafNode(item.menuId, menuTreeData.value)) {
      arr.push(Number(item.menuId))
    }
  })
  menuTreeIds.value = arr
}

init()

function nodeChange() {
  needConfirm.value = true
}
// 暴露给外层使用的切换拦截统一方法
function enterAndNext() {
  relation()
}
// 关联树 确认方法
const menuTree = ref(null)
async function relation() {
  const checkArr = menuTree.value.getCheckedNodes(false, true)
  const res = await request.post('/menu/addMenuAuthority', {
    menus: checkArr,
    authorityId: props.row.authorityId,
  })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '菜单设置成功!',
    })
  }
}

defineExpose({ enterAndNext, needConfirm })

function filterNode(value, data) {
  if (!value)
    return true
  // console.log(data.mate.title)
  return data.meta.title.includes(value)
}

watch(filterText, (val) => {
  menuTree.value.filter(val)
})
</script>

<template>
  <div>
    <div class="sticky top-0.5 z-10 bg-white">
      <el-input v-model="filterText" class="w-3/5" placeholder="筛选" />
      <el-button class="float-right" type="primary" @click="relation">
        确 定
      </el-button>
    </div>
    <div class="tree-content clear-both">
      <el-scrollbar>
        <el-tree
          ref="menuTree" :data="menuTreeData" :default-checked-keys="menuTreeIds" :props="menuDefaultProps"
          default-expand-all highlight-current node-key="ID" show-checkbox :filter-node-method="filterNode"
          @check="nodeChange"
        >
          <template #default="{ node }">
            <span class="custom-tree-node">
              <span>{{ node.label }}</span>
            </span>
          </template>
        </el-tree>
      </el-scrollbar>
    </div>
  </div>
</template>

<style lang="scss" scoped>
// 使用全局样式
</style>
