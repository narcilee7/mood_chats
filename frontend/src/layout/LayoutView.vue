<script lang="ts" setup> 
import SideBar from '../components/SideBar.vue';

</script>

<template>
 <div class="main">
  <side-bar />
  <router-view />
 </div>
</template>

<style lang='scss' scoped>
.main {
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: row;
  background-color: rgb(243, 245, 246);
}
</style>