<script setup lang="ts">

import DashboardLayout from "@/layouts/DashboardLayout.vue";
import PageHeaderCard from "@/components/shared/PageHeaderCard.vue";
import {Card, CardContent} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {onMounted, ref} from "vue";
import {LoadEntities} from "../../bindings/storyguardian/src/project/entitymanager";
import {Entity} from "../../bindings/storyguardian/src/project";
import {useToast} from "@/components/ui/toast";
import {useRouter} from "vue-router";

const {toast} = useToast();
const router = useRouter();
const entities = ref<Entity[]>([]);

onMounted(async ()=>{
  try {
    entities.value = await LoadEntities();
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
})

</script>

<template>
  <DashboardLayout>
    <PageHeaderCard>
      <p class="text-2xl leading-loose ml-2">Entities</p>
    </PageHeaderCard>

    <Card class="bg-muted/30 col-span-4 h-full">
      <CardContent class="p-6">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Created At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="entity in entities" @click="router.push('/entity/' + entity.id)" class="hover:cursor-pointer" >
              <!-- Table Data -->
              <TableCell class="font-medium">
                {{ entity.name }}
              </TableCell>
              <TableCell class="">
                {{ entity.createdAt }}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </DashboardLayout>
</template>

<style scoped>

</style>