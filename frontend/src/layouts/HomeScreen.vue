<script setup lang="ts">
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from "@/components/ui/dropdown-menu";
import {FolderOpenDot } from 'lucide-vue-next'
import {Icon} from '@iconify/vue'
import {useColorMode} from "@vueuse/core";
import {Button} from "@/components/ui/button";

const colorMode = useColorMode()
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <header class="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
      <div class="w-full flex flex-row">
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="outline">
            <Icon icon="radix-icons:moon"
                  class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"/>
            <Icon icon="radix-icons:sun"
                  class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"/>
            <span class="sr-only">Toggle theme</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <!--suppress TypeScriptValidateTypes -->
          <DropdownMenuItem @click="colorMode = 'light'">
            Light
          </DropdownMenuItem>
          <!--suppress TypeScriptValidateTypes -->
          <DropdownMenuItem @click="colorMode = 'dark'">
            Dark
          </DropdownMenuItem>
          <!--suppress TypeScriptValidateTypes -->
          <DropdownMenuItem @click="colorMode = 'system'">
            System
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
    <div class="flex flex-1 h-full w-full">
      <div class="grid md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr] w-full">
        <!--  Menu  -->
        <div class="border-r bg-muted/40 flex flex-col">
          <div class="flex-1">
            <nav class="grid items-start px-2 text-sm font-medium lg:px-4 pt-8">
              <a href="/dashboard">
                <div
                    class="flex items-center gap-3 rounded-lg px-3 py-2 text-primary transition-all hover:text-primary bg-muted">
                  <FolderOpenDot  class="h-4 w-4"/>
                  Projects
                </div>
              </a>
            </nav>
          </div>
        </div>
        <slot/>
      </div>
    </div>
  </div>
</template>
<style scoped>

</style>