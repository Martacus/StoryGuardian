<script setup lang="ts">

import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from "@/components/ui/dropdown-menu";
import {Package2} from "lucide-vue-next";
import {Icon} from "@iconify/vue";
import {Button} from "@/components/ui/button";
import {useColorMode} from "@vueuse/core";

const mode = useColorMode()
</script>

<template>
  <div class="flex min-h-screen h-full w-full flex-col">
    <header class="sticky top-0 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
      <a href="#" class="flex items-center gap-2 text-lg font-semibold md:text-base">
        <Package2 class="h-6 w-6"/>
        <span class="sr-only">Acme Inc</span>
      </a>
      <a href="/" class="text-foreground transition-colors hover:text-foreground">
        Guardian
      </a>

      <div class="ml-auto">
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
            <DropdownMenuItem @click="mode = 'light'">
              Light
            </DropdownMenuItem>
            <DropdownMenuItem @click="mode = 'dark'">
              Dark
            </DropdownMenuItem>
            <DropdownMenuItem @click="mode = 'auto'">
              System
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
    <div class="flex-grow flex items-center justify-center">
      <slot/>
    </div>
  </div>
</template>

<style scoped>

</style>