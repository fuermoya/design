<script setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref } from 'vue'
import WarningBar from '~/components/warningBar/warningBar.vue'

import Apis from '~/pages/admin/authority/components/apis.vue'
import Menus from '~/pages/admin/authority/components/menus.vue'

import request from '~/utils/request'

defineOptions({
  name: 'Authority',
})

// 表单验证函数
function mustUint(rule, value, callback) {
  if (!/^[1-9]\d*$/.test(value)) {
    return callback(new Error('请输入正整数'))
  }
  return callback()
}

// 响应式数据
const AuthorityOption = ref([
  {
    authorityId: 0,
    authorityName: '根角色',
  },
])

const drawer = ref(false)
const dialogType = ref('add')
const activeRow = ref({})

const dialogTitle = ref('新增角色')
const dialogFormVisible = ref(false)
const apiDialogFlag = ref(false)
const copyForm = ref({})

const form = ref({
  authorityId: 0,
  authorityName: '',
  parentId: 0,
})

const rules = ref({
  authorityId: [
    { required: true, message: '请输入角色ID', trigger: 'blur' },
    { validator: mustUint, trigger: 'blur', message: '必须为正整数' },
  ],
  authorityName: [
    { required: true, message: '请输入角色名', trigger: 'blur' },
  ],
  parentId: [
    { required: true, message: '请选择父角色', trigger: 'blur' },
  ],
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(999)
const tableData = ref([])
const searchInfo = ref({})

// 组件引用
const authorityForm = ref(null)
const menus = ref(null)
const apis = ref(null)

// 查询表格数据
async function getTableData() {
  const table = await request.post('/authority/getAuthorityList', {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

// 初始化数据
getTableData()

// 更新行数据
function changeRow(key, value) {
  activeRow.value[key] = value
}

// 自动进入下一个标签页
function autoEnter(activeName, oldActiveName) {
  const paneArr = [menus, apis]
  if (oldActiveName) {
    if (paneArr[oldActiveName].value.needConfirm) {
      paneArr[oldActiveName].value.enterAndNext()
      paneArr[oldActiveName].value.needConfirm = false
    }
  }
}

// 拷贝角色
function copyAuthorityFunc(row) {
  setOptions()
  dialogTitle.value = '拷贝角色'
  dialogType.value = 'copy'
  for (const k in form.value) {
    form.value[k] = row[k]
  }
  copyForm.value = row
  dialogFormVisible.value = true
}

// 打开抽屉
function opdendrawer(row) {
  drawer.value = true
  activeRow.value = row
}

// 删除角色
function deleteAuth(row) {
  ElMessageBox.confirm('此操作将永久删除该角色, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      const res = await request.post('/authority/deleteAuthority', { authorityId: row.authorityId })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!',
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消删除',
      })
    })
}

// 初始化表单
function initForm() {
  if (authorityForm.value) {
    authorityForm.value.resetFields()
  }
  form.value = {
    authorityId: 0,
    authorityName: '',
    parentId: 0,
  }
}

// 关闭窗口
function closeDialog() {
  initForm()
  dialogFormVisible.value = false
  apiDialogFlag.value = false
}

// 确定弹窗
function enterDialog() {
  authorityForm.value.validate(async (valid) => {
    if (valid) {
      form.value.authorityId = Number(form.value.authorityId)
      switch (dialogType.value) {
        case 'add': {
          const res = await request.post('/authority/createAuthority', form.value)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '添加成功!',
            })
            getTableData()
            closeDialog()
          }
        }
          break
        case 'edit': {
          const res = await request.put('/authority/updateAuthority', form.value)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '添加成功!',
            })
            getTableData()
            closeDialog()
          }
        }
          break
        case 'copy': {
          const data = {
            authority: {
              authorityId: 0,
              authorityName: '',
              datauthorityId: [],
              parentId: 0,
            },
            oldAuthorityId: 0,
          }
          data.authority.authorityId = form.value.authorityId
          data.authority.authorityName = form.value.authorityName
          data.authority.parentId = form.value.parentId
          data.authority.dataAuthorityId = copyForm.value.dataAuthorityId
          data.oldAuthorityId = copyForm.value.authorityId
          const res = await request.post('/authority/copyAuthority', data)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '复制成功！',
            })
            getTableData()
          }
        }
      }

      initForm()
      dialogFormVisible.value = false
    }
  })
}

// 设置选项
function setOptions() {
  AuthorityOption.value = [
    {
      authorityId: 0,
      authorityName: '根角色',
    },
  ]
  setAuthorityOptions(tableData.value, AuthorityOption.value, false)
}

// 设置权限选项
function setAuthorityOptions(AuthorityData, optionsData, disabled) {
  form.value.authorityId = String(form.value.authorityId)
  AuthorityData
  && AuthorityData.forEach((item) => {
    if (item.children && item.children.length) {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName,
        disabled: disabled || item.authorityId === form.value.authorityId,
        children: [],
      }
      setAuthorityOptions(
        item.children,
        option.children,
        disabled || item.authorityId === form.value.authorityId,
      )
      optionsData.push(option)
    }
    else {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName,
        disabled: disabled || item.authorityId === form.value.authorityId,
      }
      optionsData.push(option)
    }
  })
}

// 增加角色
function addAuthority(parentId) {
  initForm()
  dialogTitle.value = '新增角色'
  dialogType.value = 'add'
  form.value.parentId = parentId
  setOptions()
  dialogFormVisible.value = true
}

// 编辑角色
function editAuthority(row) {
  setOptions()
  dialogTitle.value = '编辑角色'
  dialogType.value = 'edit'
  for (const key in form.value) {
    form.value[key] = row[key]
  }
  setOptions()
  dialogFormVisible.value = true
}
</script>

<template>
  <div class="authority">
    <WarningBar title="注：角色菜单、角色api需要分开设置" />

    <div class="table-container">
      <div class="table-header">
        <h3 class="table-title">
          角色管理
        </h3>
        <div class="table-actions">
          <el-button type="primary" icon="plus" @click="addAuthority(0)">
            新增角色
          </el-button>
        </div>
      </div>

      <el-table
        :data="tableData"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        row-key="authorityId"
        style="width: 100%"
      >
        <el-table-column label="角色ID" min-width="180" prop="authorityId" />
        <el-table-column align="left" label="角色名称" min-width="180" prop="authorityName" />
        <el-table-column align="left" label="操作" width="460" class-name="table-actions-column">
          <template #default="scope">
            <div class="action-buttons">
              <el-button icon="setting" type="primary" link @click="opdendrawer(scope.row)">
                设置权限
              </el-button>
              <el-button icon="plus" type="primary" link @click="addAuthority(scope.row.authorityId)">
                新增子角色
              </el-button>
              <el-button icon="copy-document" type="primary" link @click="copyAuthorityFunc(scope.row)">
                拷贝
              </el-button>
              <el-button icon="edit" type="primary" link @click="editAuthority(scope.row)">
                编辑
              </el-button>
              <el-button icon="delete" type="primary" link @click="deleteAuth(scope.row)">
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增角色弹窗 -->
    <el-dialog v-model="dialogFormVisible" :title="dialogTitle">
      <el-form ref="authorityForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="父级角色" prop="parentId">
          <el-cascader
            v-model="form.parentId"
            style="width:100%"
            :disabled="dialogType === 'add'"
            :options="AuthorityOption"
            :props="{
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false,
            }"
            :show-all-levels="false"
            filterable
          />
        </el-form-item>
        <el-form-item label="角色ID" prop="authorityId">
          <el-input
            v-model="form.authorityId"
            :disabled="dialogType === 'edit'"
            autocomplete="off"
            maxlength="15"
          />
        </el-form-item>
        <el-form-item label="角色姓名" prop="authorityName">
          <el-input v-model="form.authorityName" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">
            取 消
          </el-button>
          <el-button type="primary" @click="enterDialog">
            确 定
          </el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer
      v-if="drawer"
      v-model="drawer"
      custom-class="auth-drawer"
      :with-header="false"
      size="40%"
      title="角色配置"
    >
      <el-tabs :before-leave="autoEnter" type="border-card">
        <el-tab-pane label="角色菜单">
          <Menus ref="menus" :row="activeRow" @change-row="changeRow" />
        </el-tab-pane>
        <el-tab-pane label="角色api">
          <Apis ref="apis" :row="activeRow" @change-row="changeRow" />
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
  </div>
</template>

<style lang="scss" scoped>
.authority {
  .el-input-number {
    margin-left: 15px;

    span {
      display: none;
    }
  }
}
</style>
