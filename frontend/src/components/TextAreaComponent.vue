<script setup lang='ts'>
import type { TextAreaType } from "@/types/index"

defineProps<{
  value: string
  type: TextAreaType
}>()

const emit = defineEmits<{
  (e: "update:value", value: string): void
}>()

const handlePublish = (e: KeyboardEvent) => {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault()
    emit("update:value", "")
  }
}


</script>

<template>
  <div class="container">
    <textarea
      placeholder="输入你现在的想法"
      class="container"
      :class="type"
      @keydown.enter="handlePublish"
    />
  </div>
</template>

<style lang='scss' scoped>
.container {
  width: 100%;
  height: 100%;

  textarea {
    border: none;
    outline: none;
    resize: none;
    padding: 1rem;
    font-size: 0.9rem;
    line-height: 1.5;
    background-color: transparent;
    color: black;
    box-shadow: 0 0 0.3rem rgba(0, 0, 0, 0.1);
    border-radius: 10px;

    scroll-behavior: smooth;
    scrollbar-3dlight-color: rgb(54, 126, 189);
    overflow-y: auto;
  }
}
</style>