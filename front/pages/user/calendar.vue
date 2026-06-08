<template>
  <div>
    <Header :user="currentUser"/>

    <div class="p-4 space-y-4">
      <UFormGroup label="包含内容" name="contentContains" :ui="{label:{base:'font-bold'}}">
        <UInput v-model="state.contentContains" placeholder="搜索内容"/>
      </UFormGroup>

      <UFormGroup label="日期范围" name="dateRange" :ui="{label:{base:'font-bold'}}">
        <UPopover :popper="{ placement: 'bottom-start' }">
          <UButton icon="i-heroicons-calendar-days-20-solid" color="white" variant="solid" class="w-full">
            {{ dateRangeLabel }}
          </UButton>

          <template #panel="{ close }">
            <div class="flex flex-col items-center sm:divide-x divide-gray-200 dark:divide-gray-800">
              <div class="flex w-full flex-col py-4 sm:flex-row">
                <UButton
                    label="所有时间"
                    color="gray"
                    variant="ghost"
                    class="rounded-none px-6"
                    :class="[!rangeEnabled ? 'bg-gray-100 dark:bg-gray-800' : 'hover:bg-gray-50 dark:hover:bg-gray-800/50']"
                    truncate
                    @click="selectAllTime(close)"
                />
                <UButton
                    v-for="(range, index) in ranges"
                    :key="index"
                    :label="range.label"
                    color="gray"
                    variant="ghost"
                    class="rounded-none px-6"
                    :class="[isRangeSelected(range.duration) ? 'bg-gray-100 dark:bg-gray-800' : 'hover:bg-gray-50 dark:hover:bg-gray-800/50']"
                    truncate
                    @click="selectRange(range.duration)"
                />
              </div>
              <DatePicker
                  :model-value="state.range"
                  @update:model-value="selectCustomRange"
                  @close="close"
              />
            </div>
          </template>
        </UPopover>
      </UFormGroup>

      <UFormGroup label="包含标签" name="tagContains" :ui="{label:{base:'font-bold'}}">
        <USelectMenu multiple v-model="state.tags" searchable :options="tags">
          <template #label>
            <span v-if="state.tags.length" class="truncate">{{ state.tags.join(', ') }}</span>
            <span v-else>选择标签</span>
          </template>
        </USelectMenu>
      </UFormGroup>
      <UFormGroup label="可见性" name="showType" :ui="{label:{base:'font-bold'}}">
        <USelectMenu v-model="state.showType"
                     :options="[{value:-1,label:'所有的'},{value:1,label:'公开的'},{value:0,label:'自己可见'}]"
                     option-attribute="label" value-attribute="value"/>
      </UFormGroup>
      <UButton class="my-2" @click="reload">搜索</UButton>
    </div>

    <div class="flex flex-col divide-y divide-[#C0BEBF]/20 ">
      <Memo v-bind:memo="m" v-for="m in memos" :key="m.id"/>
    </div>
    <div ref="loadMoreEle" class="text-xs text-center text-gray-500 py-2" @click="loadMore" v-if="hasNext">
      点击加载更多
    </div>
    <div class="text-xs text-center text-gray-500 py-2" @click="loadMore" v-else>
      已经到底啦
    </div>
  </div>
</template>

<script setup lang="ts">
import type {MemoVO, UserVO} from "~/types";
import {format, isSameDay, sub} from "date-fns";
import type {Duration} from "date-fns";
import Memo from "~/components/Memo.vue";
import {memoChangedEvent, memoReloadEvent} from "~/event";
import {useElementVisibility} from '@vueuse/core'

type DateRange = {
  start: Date
  end: Date
}

const ranges = [
  {label: '一周内', duration: {days: 7}},
  {label: '一月内', duration: {days: 31}},
  {label: '三月内', duration: {days: 90}},
  {label: '本年', duration: {months: 12}},
  {label: '近三年', duration: {years: 3}},
]
const tags = ref<string[]>([])
const currentUser = useState<UserVO>('userinfo')
const rangeEnabled = ref(false)
const state = reactive({
  page: 1,
  size: 10,
  contentContains: "",
  tags: [] as string[],
  showType: -1,
  range: {
    start: sub(new Date(), {days: 31}),
    end: new Date()
  }
})
const dateRangeLabel = computed(() => {
  if (!rangeEnabled.value) {
    return '所有时间'
  }

  return `从 ${format(state.range.start, 'yyyy-MM-dd')} 到 ${format(state.range.end, 'yyyy-MM-dd')}`
})

function isRangeSelected(duration: Duration) {
  return rangeEnabled.value && isSameDay(state.range.start, sub(new Date(), duration)) && isSameDay(state.range.end, new Date())
}

function selectAllTime(close?: () => void) {
  rangeEnabled.value = false
  close?.()
}

function selectRange(duration: Duration) {
  rangeEnabled.value = true
  state.range = {start: sub(new Date(), duration), end: new Date()}
}

function selectCustomRange(range: DateRange) {
  rangeEnabled.value = true
  state.range = range
}

const buildSearchParams = () => {
  const params: Record<string, unknown> = {
    page: state.page,
    size: state.size,
    contentContains: state.contentContains,
    showType: state.showType,
    tag: state.tags.join(','),
  }

  if (rangeEnabled.value) {
    params.start = state.range.start
    params.end = state.range.end
  }

  return params
}

const loadTags = async () => {
  const res = await useMyFetch<{
    tags: string[]
  }>("/tag/list")
  tags.value = res.tags
}


const loadMoreEle = ref(null)
const targetIsVisible = useElementVisibility(loadMoreEle)
watch(targetIsVisible, async (visible) => {
  if (visible) {
    await loadMore()
  }
})
const hasNext = ref(false)


const memos = ref<Array<MemoVO>>([])
onMounted(async () => {
  await loadTags()
  await reload()
})

const reload = async () => {
  state.page = 1
  const res = await useMyFetch<{
    list: Array<MemoVO>,
    total: number,
    hasNext: boolean
  }>('/memo/list', buildSearchParams())
  memos.value = res.list
  hasNext.value = res.hasNext
}

const loadMore = async () => {
  state.page = state.page + 1
  const res = await useMyFetch<{
    list: Array<MemoVO>,
    total: number,
    hasNext: boolean
  }>('/memo/list', buildSearchParams())
  memos.value = [...memos.value, ...res.list]
  hasNext.value = res.hasNext
}

memoReloadEvent.on(async () => {
  await reload()
})

memoChangedEvent.on(async (id: number) => {
  const res = await useMyFetch<MemoVO>('/memo/get?latest=1&id=' + id)
  const index = memos.value.findIndex(r => r.id === id)
  if (index >= 0) {
    memos.value[index] = res
  }
})

</script>

<style scoped>

</style>
