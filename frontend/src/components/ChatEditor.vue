<script setup lang="ts">
import { ref } from "vue"
import { ArrowLongUpIcon } from "@heroicons/vue/24/outline"

const emit = defineEmits<{
  (e: "update:value", value: string): void
}>()

const editorRef = ref<HTMLDivElement | null>(null)

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault()

    const text = editorRef.value?.innerText.trim() ?? ""
    if (text.length > 0) {
      emit("update:value", text)
      editorRef.value!.innerText = ""
    }
  }
}

const handlePublish = () => {
  if (editorRef.value?.innerText.trim()) {
    emit("update:value", editorRef.value.innerText.trim())
    editorRef.value.innerText = ""
  }
}

</script>

<template>
  <div class="chat-editor">
    <div class="chat-input">
      <div
        class="chat-input-box"
        ref="editorRef"
        contenteditable="true"
        @keydown="handleKeyDown"
        spellcheck="false"
        placeholder="请输入你的想法. . ."
      ></div>
    </div>
    <div class="chat-editor-action">
      <button class="publish" @click="handlePublish">
        <ArrowLongUpIcon />
      </button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.chat-editor {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

.chat-input {
  padding: 1rem;
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
  padding: 10px 0;

  &::-webkit-scrollbar {
    display: none;
  }


  .chat-input-box {
    min-height: 80px;
    max-height: calc(100% - 50px);
    width: 100%;
    outline: none;
    border: none;
    padding: 1rem;
    font-size: 1rem;
    line-height: 1.5;
    background-color: transparent;
    color: #333;
    border-radius: 10px;
    white-space: pre-wrap;

    &:empty:before {
      content: attr(placeholder);
      color: #aaa;
    }
  }
}

.chat-editor-action {
  height: 50px;
  width: 100%;
  position: relative;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: center;

  .publish {
    right: 15px;
    bottom: 15px;
    position: absolute;
    cursor: pointer;
    background-color: rgb(208, 208, 208);
    border-radius: 10px;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease-in-out;
    padding: 5px;
    font-size: 14px;
    color: white;
    font-weight: 700;

    &:hover {
      background-color: rgb(162, 162, 162);
    }
  }
}
</style>
