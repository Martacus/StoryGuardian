<script setup lang="ts">
import HomeScreen from "@/layouts/HomeScreen.vue";
import {ref} from "vue";
import {project} from "../../wailsjs/go/models";
import ProjectConfig = project.ProjectConfig;
import {Card, CardContent } from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {CreateProject, GetConfig} from "../../wailsjs/go/project/ApplicationManager";

const config = ref<ProjectConfig>();

GetConfig().then((x) => {
  config.value = x
})

function createProject(){
  CreateProject()
}
</script>

<template>
  <HomeScreen>
    <main class="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6 w-full h-full">
      <Card  class="m-auto max-w-lg lg:min-w-96">
        <CardContent class="mt-6">
          <div v-if="config" class="flex gap-2 justify-end">
            <Button @click="createProject">
              New project
            </Button>
            <Button>
              Open
            </Button>
          </div>
        </CardContent>
      </Card>
    </main>
  </HomeScreen>
</template>

<style scoped>

</style>