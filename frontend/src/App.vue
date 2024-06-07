<script lang="ts" setup>
import {useColorMode} from "@vueuse/core";
import HomeScreen from "@/layouts/HomeScreen.vue";
import {GetConfig} from "../wailsjs/go/project/ProjectManager";
import {ref} from "vue";
import {project} from "../wailsjs/go/models";
import ProjectConfig = project.ProjectConfig;

const mode = useColorMode();
const config = ref();

GetConfig().then((x)=>{
  config.value = x
})

</script>

<template>
  <HomeScreen>
    <p v-if="Object.keys(config.value.projects).length < 1">sss</p>
  </HomeScreen>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
