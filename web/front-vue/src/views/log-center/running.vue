<template>
  <NavigationBar />
  <div class="operation-container">
    <a-card title="操作日志" :bordered="false">
      <!-- 搜索表单 -->
      <a-form layout="inline" :model="searchForm" class="search-form">
        <a-form-item label="时间范围">
          <a-range-picker
            v-model:value="searchForm.timeRange"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            :placeholder="['开始时间', '结束时间']"
            @change="onTimeRangeChange"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">搜索</a-button>
            <a-button @click="resetSearch">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>

      <!-- 表格 -->
      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        :scroll="{ x: 1300 }"
        class="operation-table"
      >
        <!-- 通用 bodyCell 插槽：保持你现有的渲染规则，并加入 action 列渲染 -->
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'status'">
            <a-tag :color="text === 200 ? 'success' : 'error'">
              {{ text }}
            </a-tag>
          </template>

          <template v-else-if="column.dataIndex === 'startTime'">
            {{ formatDate(text) }}
          </template>

          <!-- 新增：操作列渲染 -->
          <template v-else-if="column.dataIndex === 'action'">
            <a-button size="small" @click="openDetail(record.path)">查看</a-button>
          </template>

          <!-- 默认回退显示 -->
          <template v-else>
            {{ text }}
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 详情 Modal -->
    <a-modal
      title="日志详情"
      :visible="detailModalVisible"
      @cancel="closeModal"
      @ok="closeModal"
      width="900px"
      :footer="null"
      destroyOnClose
    >
      <a-spin :spinning="detailLoading">
        <div v-if="detailLines.length">
          <div class="log-list">
            <div
              class="log-line"
              v-for="(line, idx) in detailLines"
              :key="idx"
            >
              <a-tag :color="line.level === 'ERROR' ? 'red' : 'blue'">
                {{ line.level || 'RAW' }}
              </a-tag>
              <span class="time">{{ line.time || '' }}</span>
              <span class="msg">{{ line.msg || line.raw }}</span>
            </div>
          </div>
        </div>

        <div v-else class="text-center">
          暂无数据
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import NavigationBar from '@/components/NavigationBar.vue'
import { useTable } from '@/tools/page/common'
import { logByRunningColumns } from '@/tools/page/columns'
import { formatDate, newTimeRangeHandler } from '@/tools/page/time'
import { logByRunningOne } from '@/tools/api'

/*----------------------------------------全局变量----------------------------------------*/
const {
  tableData,
  loading,
  pagination,
  searchForm,
  fetchData,
  resetSearch,
  handleSearch,
  handleTableChange,
} = useTable('logs-running')
const onTimeRangeChange = newTimeRangeHandler(searchForm)

const columns = computed(() => {
  // 展开旧列，然后追加操作列（固定右侧）
  return [
    ...logByRunningColumns,
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      width: 30,
      align: 'center',
    },
  ]
})

/*----------------------------------------详情 Modal 相关----------------------------------------*/
const detailModalVisible = ref(false)
const detailLoading = ref(false)
// detailData 保存后端响应的 body（通常是字符串）
const detailData = ref(null)

// 把 detailData（可能是 { data: "..." } 或直接字符串）按换行切分，并尝试解析 JSON
const detailLines = computed(() => {
  if (!detailData.value) return []
  // rawText 尽量从几种可能的结构里取到字符串
  const rawText =
    typeof detailData.value === 'string'
      ? detailData.value
      : detailData.value.data !== undefined
      ? detailData.value.data
      : detailData.value.toString()

  // 按 \n 切分，过滤空行
  const lines = rawText.split(/\r?\n/).filter((l) => l && l.trim() !== '')

  // 将每行尽量解析成对象 { level, time, msg }，否则保留原始 raw 字段
  return lines.map((l) => {
    try {
      const obj = JSON.parse(l)
      return {
        level: obj.level,
        time: obj.time,
        msg: obj.msg,
        raw: l,
      }
    } catch (e) {
      return { raw: l }
    }
  })
})

async function openDetail(path) {
  if (!path) {
    message.warn('当前行没有 path 可用')
    return
  }
  detailModalVisible.value = true
  detailLoading.value = true
  detailData.value = null

  try {
    // 你的 API 实际签名我不清楚，之前你用过 logByRunningOne(path) 所以这里保持相同调用
    const res = await logByRunningOne(path)
    // 让 detailData 成为字符串或后端的数据体（优先取 res.data）
    detailData.value =
      res && res.data !== undefined && typeof res.data === 'string'
        ? res.data
        : res && res.data !== undefined
        ? res.data
        : typeof res === 'string'
        ? res
        : res
  } catch (err) {
    console.error(err)
    message.error('获取日志详情失败：' + (err?.message || err))
  } finally {
    detailLoading.value = false
  }
}

function closeModal() {
  detailModalVisible.value = false
  detailData.value = null
}

/*----------------------------------------生命周期----------------------------------------*/
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 根据需要加样式 */
.operation-container {
  padding: 16px;
}
.operation-table ::v-deep .ant-table-cell {
  vertical-align: middle;
}

/* 日志列表样式 */
.log-list {
  max-height: 60vh;
  overflow: auto;
  padding-right: 8px;
}
.log-line {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 8px;
  border-bottom: 1px solid #f0f0f0;
  word-break: break-word;
}
.log-line .time {
  color: #888;
  font-size: 12px;
  min-width: 220px;
}
.log-line .msg {
  white-space: pre-wrap;
  flex: 1;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, "Roboto Mono", "Courier New", monospace;
  font-size: 13px;
}
.text-center {
  text-align: center;
  padding: 16px 0;
}
</style>
