<script setup lang="ts">
import HomeScreen from "@/layouts/HomeScreen.vue";
import {ref} from "vue";
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {CreateProject, GetConfig} from "../../bindings/storyguardian/project/applicationmanager";
import {ApplicationConfig} from "../../bindings/storyguardian/project";
import {Separator} from "@/components/ui/separator";

const config = ref<ApplicationConfig>();

GetConfig().then((x) => {
  config.value = x
})

function createProject(){
  CreateProject()
}

// function openProject(){
//   OpenProject()
// }
</script>

<template>
  <HomeScreen>
    <main class="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6 w-full h-full">
      <Card  class="m-auto max-w-lg lg:min-w-96">
        <CardHeader>
          <div class="flex flex-row gap-2 items-center">
            <h1 class="text-2xl">Projects</h1>
            <div v-if="config" class="flex gap-2 justify-end w-full">
              <Button @click="createProject">
                New project
              </Button>
              <Button>
                Open
              </Button>
            </div>
          </div>
          <Separator class="mt-4"/>
        </CardHeader>
        <CardContent class="flex flex-col" v-if="config">
          <div v-for="item in config.projects" class="hover:bg-white mt-2">
            <h1>{{item.name}}</h1>
            <p class="text-sm text-gray-500">{{item.location}}</p>
          </div>
        </CardContent>
      </Card>
    </main>
  </HomeScreen>
</template>

<style scoped>

</style>