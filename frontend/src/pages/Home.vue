<script setup lang="ts">
import HomeScreen from "@/layouts/HomeScreen.vue";
import {onMounted, ref} from "vue";
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {Separator} from "@/components/ui/separator";
import {useRouter} from "vue-router";
import {useToast} from "@/components/ui/toast";
import {ApplicationConfig} from "../../bindings/storyguardian/src/project";
import {CreateProject, GetConfig, OpenProject} from "../../bindings/storyguardian/src/project/applicationmanager";
import {NewStory} from "../../bindings/storyguardian/src/project/storymanager";

const {toast} = useToast()
const config = ref<ApplicationConfig>();
const router = useRouter()

onMounted(async () => {
  config.value = await GetConfig()
});


function createProject(){
  CreateProject().then((projectDirectory: string) => {
    NewStory(projectDirectory).then(story => {
      router.push('/dashboard/' + story?.id)
    });
  });
}

async function openProject(id: string){
  try{
    await OpenProject(id);
    await router.push('/dashboard/' + id)
  } catch(error: any){
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error,
    });
  }
}
</script>

<template>
  <HomeScreen>
    <Card class="m-auto max-w-lg lg:min-w-96 col-span-4">
      <CardHeader>
        <div class="flex flex-row gap-2 items-center">
          <h1 class="text-2xl">Stories</h1>
          <div v-if="config" class="flex gap-2 ml-4 justify-end w-full">
            <Button @click="createProject">
              New story
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