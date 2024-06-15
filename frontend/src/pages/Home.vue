<script setup lang="ts">
import HomeScreen from "@/layouts/HomeScreen.vue";
import {ref} from "vue";
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {CreateProject, GetConfig} from "../../bindings/storyguardian/project/applicationmanager";
import {ApplicationConfig} from "../../bindings/storyguardian/project";
import {Separator} from "@/components/ui/separator";
import {useRouter} from "vue-router";

const config = ref<ApplicationConfig>();
const router = useRouter()

GetConfig().then((x) => {
  config.value = x
})

function createProject(){
  CreateProject().then((id) => {
    router.push('/dashboard/' + id)
  })
}

function openProject(id: string){
  router.push('/dashboard/' + id)
}
</script>

<template>
  <HomeScreen>
    <Card  class="m-auto max-w-lg lg:min-w-96">
      <CardHeader>
        <div class="flex flex-row gap-2 items-center">
          <h1 class="text-2xl">Projects</h1>
          <div v-if="config" class="flex gap-2 ml-4 justify-end w-full">
            <Button @click="createProject">
              New project
            </Button>
            <Button>
              Open
            </Button>
            <Button>
              Import
            </Button>
          </div>
        </div>
        <Separator class="mt-4"/>
      </CardHeader>
      <CardContent class="flex flex-col" v-if="config">
        <div v-for="item in config.projects" class="hover:bg-muted hover:cursor-pointer rounded p-2" @click="openProject(item.id)">
          <h1 class="text-white">{{item.name}}</h1>
          <p class="text-sm text-gray-500">{{item.location}}</p>
          <p class="text-sm text-gray-500">{{item.id}}</p>
        </div>
      </CardContent>
    </Card>
  </HomeScreen>
</template>

<style scoped>

</style>