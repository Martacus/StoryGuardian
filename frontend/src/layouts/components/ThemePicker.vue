<script setup lang="ts">

import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from "@/components/ui/dropdown-menu";
import {Icon} from "@iconify/vue";
import {Button} from "@/components/ui/button";
import {useColorMode} from "@vueuse/core";

const colorMode = useColorMode()
</script>

<template>
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
</template>

<style scoped>

</style>