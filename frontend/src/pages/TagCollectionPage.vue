<script setup lang="ts">

import {useToast} from "@/components/ui/toast";
import {onMounted, ref} from "vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import PageHeaderCard from "@/components/shared/PageHeaderCard.vue";
import {Card, CardContent} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {GetStoryTags} from "../../bindings/storyguardian/src/project/storymanager";

const {toast} = useToast()
const tags = ref<string[]>([])

onMounted(async ()=>{
  try {
    tags.value = await GetStoryTags();
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
      <p class="text-2xl leading-loose ml-2">Tags</p>
    </PageHeaderCard>

    <Card class="bg-muted/30 col-span-4 h-full">
      <CardContent class="p-6">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="tag in tags" @click="" class="hover:cursor-pointer">
              <!-- Table Data -->
              <TableCell class="font-medium">
                {{ tag }}
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