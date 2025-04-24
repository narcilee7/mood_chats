<script setup lang='ts'>
import { ChevronUpIcon, ChevronDownIcon  } from "@heroicons/vue/20/solid";
import { ChatBubbleLeftRightIcon, ClockIcon } from "@heroicons/vue/24/outline";
import { ref, onMounted, onUnmounted } from "vue";

const footChevronUp = ref(false);

const handleFooterClick = () => {
  footChevronUp.value = !footChevronUp.value;
}

const handleCreateNewChat = () => {
  console.log('handleCreateNewChat')
}

// 监听键盘事件
const handleKeyDown = (e: KeyboardEvent) => {
  try {
    if (e.ctrlKey && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      handleCreateNewChat()
    }
  } catch (error) {
    console.error('Error handling keyboard event:', error)
  }
}

onMounted(() => {
  try {
    window.addEventListener('keydown', handleKeyDown)
  } catch (error) {
    console.error('Error adding event listener:', error)
  }
})

onUnmounted(() => {
  try {
    window.removeEventListener('keydown', handleKeyDown)
  } catch (error) {
    console.error('Error removing event listener:', error)
  }
})
</script>

<template>
 <div class="sidebar">
  <div class="header side-item">
    <h2>Mood-ChatBot</h2>
  </div>
  <div class="body">
    <button class="body_btn">
      <div class="body_btn_left" @click="handleCreateNewChat">
        <ChatBubbleLeftRightIcon class="w-8 h-8 text-gray-400" />
        <span>开启新对话</span>
      </div>
      <div class="body_btn_right">
        <span>Ctrl</span>
        <span>K</span>
      </div>
    </button>
    <div class="body_history">
      <div class="body_history_title">
        <ClockIcon class="w-8 h-8 text-gray-400" />
        <h2>历史对话</h2>
      </div>
    </div>
  </div>
  <div class="footer side-item">
    <div class="info">
      <img src="../assets/img/icon.jpg" alt="user icon">
      <span>cc</span>
    </div>
    <div class="action" @click="handleFooterClick">
      <ChevronUpIcon class="w-10 h-10 ext-gray-400 cursor-pointer" v-if="footChevronUp" />
      <ChevronDownIcon class="w-10 h-10 ext-gray-400 cursor-pointer" v-else />
    </div>
  </div>
 </div>
</template>

<style lang='scss' scoped>
.sidebar {
  width: 300px;
  min-width: 300px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  background-color: rgb(243, 245, 246);
}

.side-item {
  width: 100%;
  height: 50px;
  min-height: 50px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 15px;

  h2 {
    font-size: 20px;
    color: #000;
    font-weight: 600;
    letter-spacing: 0.5px;
    font-family: 'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;
  }
}

.body {
  width: 100%;
  height: calc(100% - 100px);
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  border-radius: 20px;
  padding: 0 5px;

  &_btn {
    width: 100%;
    height: 50px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    margin-top: 5px;
    border-radius: 10px;
    background-color: #fff;
    border: 1px solid rgb(235, 230, 230);

    &:hover {
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }


    &_left {
      display: flex;
      gap: 4px;
      align-items: center;
      cursor: pointer;
    }

    &_right {
      cursor: pointer;
      display: flex;
      gap: 2px;
      
      span {
        border-radius: 5px;
        background-color: rgb(243, 245, 246);
        background-color: rgb(242, 242, 242);
        padding: 2px;
      }
    }
  }

  &_history {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;

    &_title {
      width: 100%;
      padding: 10px;
      display: flex;
      justify-content: flex-start;
      align-items: center;
      gap: 10px;
    }
  }
}

.footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 15px;

  .info {
    height: 100%;
    display: flex;
    align-items: center;
    gap: 10px;
    
    img {
      object-fit: contain;
      width: 40px;
      height: 40px;
      border-radius: 50%;
    }

    span {
      font-size: 20px;
    }
  }

  .action {
    height: 100%;
    width: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

@media screen and (max-width: 768px) {
  // 小于400直接小时
  .sidebar {
    display: none;
    // 丝滑一些
    transition: all 0.5s ease-in-out;
  }
}

</style>