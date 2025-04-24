<script setup lang='ts'>
import ChatEditor from '@/components/ChatEditor.vue';

</script>

<template>
  <div class="chat_container">
    <div class="chat_container_header"></div>
    <div class="chat_container_logo">Mood-ChatBot</div>
    <div class="chat_container_textarea">
      <chat-editor 
        type="index"
      />
    </div>
    <div class="chat_container_tags"></div>
    <div class="chat_container_footer"></div>
  </div>
</template>

<style lang='scss' scoped>
.chat_container {
  width: calc(100% - 24px);
  height: calc(100% - 24px);
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 12px 12px 12px 12px;
  border-radius: 10px;
  background-color: white;

  &_logo {
    margin-bottom: 20px;
    font-size: 30px;
    font-weight: 600;
    letter-spacing: 3px;
    font-family: 'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;
  }

  &_header {
    width: 100%;
    height: 250px;
  }

  &_footer {
    width: 100%;
    height: 300px;
  }

  &_textarea {
    width: 80%;
    height: 230px;
  }
}
</style>