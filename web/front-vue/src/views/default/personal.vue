<template>
  <NavigationBar />
  <div class="profile-container">
    <a-card :bordered="false" :bodyStyle="{ padding: '24px' }">
      <div class="profile-header">
        <h2>个人信息</h2>
      </div>

      <div class="profile-content">
        <!-- 用户头像 -->
        <div class="avatar-section">
          <a-avatar :size="80" :src="user.avatarUrl" style="background-color: #00a1d6; fontSize: 24px;">
            {{ getAvatarText }}
          </a-avatar>
          <div class="avatar-text">
            <h3>{{ user.name }}</h3>
            <p>账号：{{ user.account }}</p>
          </div>
        </div>

        <!-- 信息卡片 -->
        <a-card :bordered="true" class="info-card">
          <a-descriptions title="基本信息" :column="1" size="middle" bordered>
            <a-descriptions-item label="用户名">
              <a-tag color="blue">{{ user.name }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="邮箱">
              <span>{{ user.email }}</span>
            </a-descriptions-item>
            <a-descriptions-item label="账号">
              {{ user.account }}
            </a-descriptions-item>
            <a-descriptions-item label="角色">
              <a-tag v-for="role in user.roles" :key="role.id" color="geekblue">
                {{ role.name }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="创建时间">
              {{ formatDate(user.created_at) }}
            </a-descriptions-item>
            <a-descriptions-item label="更新时间">
              {{ user.updated_at === '0001-01-01T00:00:00Z' ? '暂无更新' : formatDate(user.updated_at) }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- Token 信息（可选，调试用） -->
        <!-- <a-card :bordered="true" class="info-card" style="margin-top: 24px;">
          <a-descriptions title="Token 信息" :column="1" size="small" bordered>
            <a-descriptions-item label="Access Token">
              <a-tooltip :title="data.access_token">
                <span>{{ data.access_token.substring(0, 50) }}...</span>
              </a-tooltip>
            </a-descriptions-item>
            <a-descriptions-item label="Refresh Token">
              <a-tooltip :title="data.refresh_token">
                <span>{{ data.refresh_token.substring(0, 50) }}...</span>
              </a-tooltip>
            </a-descriptions-item>
          </a-descriptions>
        </a-card> -->
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import NavigationBar from '@/components/NavigationBar.vue'
import {formatDate}from '@/tools/page/time'
// 响应式数据
const data = ref(null)
const user = ref({})

// 获取用户信息
const fetchUserInfo = () => {
  // 模拟从登录接口返回的数据（实际项目中应从 API 获取）
  const mockResponse = {
    code: 2000,
    message: {
      'zh-CN': '请求成功',
      'en-US': 'success'
    },
    data: {
      user: {
        id: 1,
        created_at: '2025-09-19T16:51:02Z',
        updated_at: '0001-01-01T00:00:00Z',
        name: 'chenyouwei',
        email: 'chenyouwei3@outlook.com',
        account: '21480',
        password: '$2a$12$nHMpTEE/xsbWmO3wtbdAluCzkh5kWBpWn.P60pjeteMiEv.fI8246',
        avatarUrl: '',
        roles: [
          {
            id: 1,
            created_at: '0001-01-01T00:00:00Z',
            updated_at: '0001-01-01T00:00:00Z',
            name: 'super admin',
            desc: ''
          }
        ]
      },
      access_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiY2hlbnlvdXdlaSIsImV4cCI6MTc1ODM3OTEyMiwiaXNzIjoiZ2luLXdlYiJ9.jKOfwFWSK7vQLpej3IRY3tMq6OsIVTFD9jD0CUl27dQ',
      refresh_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTg0NTgzMjIsImlzcyI6Imdpbi13ZWIifQ.aMWyDjiMcS_oeY_oQlwx0MZ4kUY9iO1lI98jQnLMQeM'
    }
  }

  data.value = mockResponse.data
  user.value = mockResponse.data.user
}

// 计算头像显示文字
const getAvatarText = user.value?.name
  ? user.value.name.charAt(0).toUpperCase()
  : user.value?.account?.charAt(0).toUpperCase() || 'U'


// 页面加载时获取数据
onMounted(() => {
  fetchUserInfo()
})
</script>

<style lang="less" scoped>

.profile-header h2 {
  margin: 0;
  font-size: 18px;
  color: #18191c;
  font-weight: 600;
}

.profile-content {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
}

.avatar-section .avatar-text h3 {
  margin: 0;
  font-size: 16px;
  color: #18191c;
  font-weight: 500;
}

.avatar-section .avatar-text p {
  margin: 4px 0 0;
  color: #9499a0;
  font-size: 14px;
}

.info-card :deep(.ant-descriptions-title) {
  font-size: 16px;
  color: #18191c;
  margin-bottom: 16px;
}

.info-card :deep(.ant-tag) {
  border-radius: 4px;
}
</style>