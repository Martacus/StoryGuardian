<script setup lang="ts">
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from "@/components/ui/dropdown-menu";
import {Bell, CircleUser, Home, Menu, Package2, Plus, ShoppingCart} from 'lucide-vue-next'
import {Icon} from '@iconify/vue'
import {useColorMode} from "@vueuse/core";
import {Button} from "@/components/ui/button";
import {Sheet, SheetContent, SheetTrigger} from "@/components/ui/sheet";

const colorMode = useColorMode()
</script>

<template>
  <div class="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
    <div class="hidden border-r bg-muted/40 md:block">
      <div class="flex h-full max-h-screen flex-col gap-2">
        <div class="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
          <a href="/" class="flex items-center gap-2 font-semibold">
            <Package2 class="h-6 w-6"/>
            <span class="">Guardian</span>
          </a>
          <Button variant="outline" size="icon" class="ml-auto h-8 w-8">
            <Bell class="h-4 w-4"/>
            <span class="sr-only">Toggle notifications</span>
          </Button>
        </div>
        <div class="flex-1">
          <nav class="grid items-start px-2 text-sm font-medium lg:px-4">
            <a href="/dashboard">
              <div
                  class="flex items-center gap-3 rounded-lg px-3 py-2 text-primary transition-all hover:text-primary bg-muted">
                <Home class="h-4 w-4"/>
                Dashboard
              </div>

            </a>
            <a href="#"
               class="flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary ">
              <ShoppingCart class="h-4 w-4"/>
              Entities
            </a>
          </nav>
        </div>
      </div>
    </div>
    <div class="flex flex-col">
      <header class="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
        <Sheet>
          <SheetTrigger as-child>
            <Button variant="outline" size="icon" class="shrink-0 md:hidden">
              <Menu class="h-5 w-5"/>
              <span class="sr-only">Toggle navigation menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" class="flex flex-col">
            <nav class="grid gap-2 text-lg font-medium">
              <a href="#" class="flex items-center gap-2 text-lg font-semibold">
                <Package2 class="h-6 w-6"/>
                <span class="sr-only">Acme Inc</span>
              </a>
              <a href="#"
                 class="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground">
                <Home class="h-5 w-5"/>
                Dashboard
              </a>
            </nav>
          </SheetContent>
        </Sheet>
        <div class="w-full flex flex-row">
          <Button variant="outline" class="ml-2">
            <Plus class="size-4" @click=""></Plus>
          </Button>
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
            <DropdownMenuItem @click="colorMode = 'light'">
              Light
            </DropdownMenuItem>
            <DropdownMenuItem @click="colorMode = 'dark'">
              Dark
            </DropdownMenuItem>
            <DropdownMenuItem @click="colorMode = 'system'">
              System
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="secondary" size="icon" class="rounded-full">
              <CircleUser class="h-5 w-5"/>
              <span class="sr-only">Toggle user menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuSeparator/>
            <DropdownMenuItem>Settings</DropdownMenuItem>
            <DropdownMenuItem>Support</DropdownMenuItem>
            <DropdownMenuSeparator/>
            <DropdownMenuItem @click="">Logout</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </header>
      <slot/>

    </div>
  </div>
</template>

<style scoped>

</style>