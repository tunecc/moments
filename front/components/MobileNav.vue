<template>
  <UModal
    v-model="open"
    :ui="{
      container:
        'fixed top-0 left-0 right-0 bottom-0 flex justify-center items-center backdrop-blur',
    }"
  >
    <div
      v-if="global.userinfo.token"
      class="flex flex-col items-center p-4 pt-8 text-gray-500 dark:text-white"
      @click="navigate('/new')"
      title="发表"
    >
      <span
        class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
      >
        <UIcon name="i-carbon-camera" class="w-7 h-7 cursor-pointer" />
      </span>
      <span class="text-sm mt-1">发表</span>
    </div>
    <div
      class="flex items-center justify-center gap-3 p-4 text-gray-500 dark:text-white min-h-[120px]"
    >
      <div class="flex flex-col items-center gap-1">
        <button
          type="button"
          class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
          :title="themeModeLabel"
          :aria-label="themeModeLabel"
          @click="toggleMode"
        >
          <UIcon :name="themeModeIcon" class="w-6 h-6 cursor-pointer" />
        </button>
        <span @click="toggleMode" class="text-sm mt-1">
          {{ themeModeText }}
        </span>
      </div>
      <div
        v-if="$route.path !== '/user/calendar' && global.userinfo.token"
        class="flex flex-col items-center"
        @click="navigate('/user/calendar')"
        title="日历检索"
      >
        <span
          class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
        >
          <UIcon name="i-jam-search-folder" class="w-6 h-6 cursor-pointer" />
        </span>
        <span class="text-sm mt-1">检索</span>
      </div>
      <div
        v-if="$route.path == '/'"
        class="flex flex-col items-center"
        @click="navigate('/friend')"
        title="友情链接"
      >
        <span
          class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
        >
          <UIcon name="i-carbon-friendship" class="w-6 h-6 cursor-pointer" />
        </span>
        <span class="text-sm mt-1">友链</span>
      </div>
      <div
        v-if="$route.path !== '/sys/settings' && global.userinfo.id === 1"
        class="flex flex-col items-center"
        @click="navigate('/sys/settings')"
        title="系统设置"
      >
        <span
          class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
        >
          <UIcon name="i-carbon-settings" class="w-6 h-6 cursor-pointer" />
        </span>
        <span class="text-sm mt-1">系统</span>
      </div>
      <div
        v-if="$route.path !== '/user/settings' && global.userinfo.token"
        class="flex flex-col items-center"
        @click="navigate('/user/settings')"
        title="用户中心"
      >
        <span
          class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
        >
          <UIcon name="i-carbon-user-avatar" class="w-6 h-6 cursor-pointer" />
        </span>
        <span class="text-sm mt-1">用户</span>
      </div>
      <div
        v-if="$route.path === '/user/settings' && global.userinfo.token"
        class="flex flex-col items-center"
        title="登出"
        @click="logout"
      >
        <span
          class="flex items-center bg-gray-200/75 dark:bg-gray-800/75 p-3 rounded-full"
        >
          <UIcon name="i-carbon-logout" class="w-6 h-6 cursor-pointer" />
        </span>
        <span class="text-sm mt-1">退出</span>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { useGlobalState } from "~/store";

const global = useGlobalState();
const mode = useColorMode();
const open = useState<boolean>("sidebarOpen", () => false);
const resolvedModeLabel = computed(() =>
  mode.value === "dark" ? "暗色" : "亮色"
);
const themeModeText = computed(() => {
  if (mode.preference === "system") {
    return `自动(${resolvedModeLabel.value})`;
  }

  return mode.preference === "dark" ? "暗色" : "亮色";
});
const themeModeLabel = computed(() => `显示模式：${themeModeText.value}`);
const themeModeIcon = computed(() => {
  if (mode.preference === "system") {
    return "i-carbon-screen";
  }

  return mode.preference === "dark" ? "i-carbon-moon" : "i-carbon-sun";
});

const toggleMode = () => {
  if (mode.preference === "system") {
    mode.preference = "light";
  } else if (mode.preference === "light") {
    mode.preference = "dark";
  } else {
    mode.preference = "system";
  }
};

const navigate = async (url: string) => {
  open.value = false;
  await navigateTo(url);
};

const logout = async () => {
  global.value.userinfo = {};
  await navigateTo("/");
};
</script>

<style scoped></style>
