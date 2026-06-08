<template>
  <div
    v-if="$route.path !== '/new' && $route.path.indexOf('/edit/') < 0"
    class="header relative mb-14"
  >
    <div
      v-if="$route.path !== '/' && $route.path.indexOf('/memo/') < 0"
      :class="{ 'bg-[#4c4c4c]/80 z-10': y > 100 }"
      class="flex fixed justify-between items-center p-4 w-full md:w-[567px] text-white top-0"
    >
      <NuxtLink class="flex items-center" title="返回主页">
        <UIcon
          @click="navigateTo('/')"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer mr-4"
        />
        <span v-if="$route.path === '/user/calendar'">日历检索</span>
        <span v-else-if="$route.path === '/sys/settings'">系统设置</span>
        <span v-else-if="$route.path === '/user/settings'">用户中心</span>
        <span v-else-if="$route.path.indexOf('/tags/') >= 0">
          {{ route.params.tag || "话题专栏" }}
        </span>
        <span v-else>
          <span v-if="!global.userinfo.token && $route.path === '/user/login'">
            登录
          </span>
          <span
            v-else-if="!global.userinfo.token && $route.path === '/user/reg'"
          >
            注册
          </span>
          <span v-else>{{ props.user.nickname }} 的空间</span>
        </span>
      </NuxtLink>
      <NuxtLink
        v-if="$route.path === '/user/settings' && global.userinfo.token"
        class="hidden sm:flex"
        title="登出"
        @click="logout"
      >
        <UIcon name="i-carbon-logout" class="w-5 h-5 cursor-pointer" />
      </NuxtLink>
    </div>

    <div
      class="dark:bg-neutral-800 hidden sm:flex sm:absolute sm:-right-10 sm:rounded sm:p-2 sm:flex-col sm:w-fit justify-end shadow w-full flex-row top-0 p-1 flex gap-2 bg-white"
    >
      <button
        type="button"
        class="flex text-[#9fc84a] cursor-pointer"
        :title="themeModeLabel"
        :aria-label="themeModeLabel"
        @click="toggleMode"
      >
        <UIcon :name="themeModeIcon" class="w-5 h-5" />
      </button>

      <NuxtLink v-if="global.userinfo.token" to="/new" title="发表">
        <UIcon
          name="i-carbon-camera"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path !== '/user/calendar' && global.userinfo.token"
        to="/user/calendar"
        title="日历检索"
      >
        <UIcon
          name="i-jam-search-folder"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path !== '/sys/settings' && global.userinfo.id === 1"
        to="/sys/settings"
        title="系统设置"
      >
        <UIcon
          name="i-carbon-settings"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path !== '/user/settings' && global.userinfo.token"
        to="/user/settings"
        title="用户中心"
      >
        <UIcon
          name="i-carbon-user-avatar"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink v-if="!global.userinfo.token" to="/user/login" title="登录">
        <UIcon
          name="i-octicon-sign-in-16"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
    </div>

    <img class="header-img w-full" :src="props.user.coverUrl" alt="" />
    <div class="absolute right-2 bottom-[-40px]">
      <div class="userinfo flex flex-col">
        <div class="flex flex-row items-center gap-4 justify-end">
          <div class="username text-lg font-bold text-white">
            {{ props.user.nickname }}
          </div>
          <img
            :src="props.user.avatarUrl"
            class="avatar w-[70px] h-[70px] rounded-xl"
          />
        </div>
        <div class="slogon text-gray truncate w-full text-end text-xs mt-2">
          {{ props.user.slogan }}
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import type { UserVO } from "~/types";
import { useGlobalState } from "~/store";

const global = useGlobalState();
const route = useRoute();

const props = defineProps<{ user: UserVO }>();
const mode = useColorMode();
const { y } = useWindowScroll();
const resolvedModeLabel = computed(() =>
  mode.value === "dark" ? "暗色" : "亮色"
);
const themePreferenceLabel = computed(() => {
  if (mode.preference === "system") {
    return `自动（当前${resolvedModeLabel.value}）`;
  }

  return mode.preference === "dark" ? "暗色" : "亮色";
});
const themeModeLabel = computed(() => `显示模式：${themePreferenceLabel.value}`);
const themeModeIcon = computed(() => {
  if (mode.preference === "system") {
    return "i-carbon-screen";
  }

  return mode.preference === "dark" ? "i-carbon-moon" : "i-carbon-sun";
});

const logout = async () => {
  global.value.userinfo = {};
  await navigateTo("/");
};

const toggleMode = () => {
  if (mode.preference === "system") {
    mode.preference = "light";
  } else if (mode.preference === "light") {
    mode.preference = "dark";
  } else {
    mode.preference = "system";
  }
};
</script>

<style scoped></style>
