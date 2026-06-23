<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Link as LinkIcon, ChevronUp, ChevronDown, ChevronRight } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { useLinkStore } from '@/stores/linkStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import type { Link } from '../../../bindings/litguardian/internal'
import type { Entity } from '../../../bindings/litguardian/internal'

const entityStore = useEntityStore()
const linkStore = useLinkStore()
const navStore = useNavigationStore()

onMounted(() => {
  if (entityStore.entities.length === 0) entityStore.loadEntities()
  linkStore.loadAllLinks()
})

// ── View state ──────────────────────────────────────────────────────────────
const mode = ref<'relation' | 'entity'>('relation')
const query = ref('')

type SortKey = 'from' | 'type' | 'to'
const sortKey = ref<SortKey>('from')
const sortDir = ref<'asc' | 'desc'>('asc')

function toggleSort(key: SortKey) {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortDir.value = 'asc'
  }
}

// ── Lookups ───────────────────────────────────────────────────────────────
const entityById = computed(
  () => new Map<string, Entity>(entityStore.entities.map(e => [e.id, e])),
)

function entityName(id: string) {
  return entityById.value.get(id)?.name ?? 'Unknown'
}

/** A link matches the search query if any entity name, the type, or the description contains it. */
function matchesQuery(link: Link) {
  const q = query.value.toLowerCase().trim()
  if (!q) return true
  return (
    entityName(link.fromEntityId).toLowerCase().includes(q) ||
    entityName(link.toEntityId).toLowerCase().includes(q) ||
    (link.type ?? '').toLowerCase().includes(q) ||
    (link.description ?? '').toLowerCase().includes(q)
  )
}

// ── By Relation rows ──────────────────────────────────────────────────────
interface RelationRow {
  link: Link
  fromName: string
  fromType: string
  toName: string
  toType: string
}

const relationRows = computed<RelationRow[]>(() => {
  const rows = linkStore.allLinks
    .filter(matchesQuery)
    .map(link => ({
      link,
      fromName: entityName(link.fromEntityId),
      fromType: entityById.value.get(link.fromEntityId)?.type ?? '',
      toName: entityName(link.toEntityId),
      toType: entityById.value.get(link.toEntityId)?.type ?? '',
    }))

  const dir = sortDir.value === 'asc' ? 1 : -1
  rows.sort((a, b) => {
    const av = sortKey.value === 'from' ? a.fromName : sortKey.value === 'to' ? a.toName : a.link.type
    const bv = sortKey.value === 'from' ? b.fromName : sortKey.value === 'to' ? b.toName : b.link.type
    return (av ?? '').localeCompare(bv ?? '', undefined, { sensitivity: 'base' }) * dir
  })
  return rows
})

// ── By Entity groups ──────────────────────────────────────────────────────
const expanded = ref<Set<string>>(new Set())

function toggleExpanded(id: string) {
  const next = new Set(expanded.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  expanded.value = next
}

interface EntityGroup {
  entity: Entity
  links: Link[]
}

const entityGroups = computed<EntityGroup[]>(() => {
  const q = query.value.toLowerCase().trim()
  const groups: EntityGroup[] = []
  for (const entity of entityStore.entities) {
    const links = linkStore.allLinks.filter(
      l => l.fromEntityId === entity.id || l.toEntityId === entity.id,
    )
    if (links.length === 0) continue
    const nameMatches = !q || entity.name.toLowerCase().includes(q)
    if (!nameMatches && !links.some(matchesQuery)) continue
    groups.push({ entity, links })
  }
  groups.sort((a, b) => a.entity.name.localeCompare(b.entity.name, undefined, { sensitivity: 'base' }))
  return groups
})

/** Name of the entity on the other side of a link, relative to a given entity. */
function otherEntityName(link: Link, entityId: string) {
  const otherId = link.fromEntityId === entityId ? link.toEntityId : link.fromEntityId
  return entityName(otherId)
}

// ── Navigation ────────────────────────────────────────────────────────────
function openRelation(link: Link) {
  navStore.navigateTo({
    view: 'links',
    entityId: link.fromEntityId,
    linkId: link.id,
    categoryName: null,
  })
}
</script>

<template>
  <div class="flex flex-col h-full">

    <!-- Header -->
    <div class="flex items-center gap-4 px-6 py-4 shrink-0">
      <h1 class="text-xl font-semibold">Relations</h1>

      <!-- Segmented toggle -->
      <div class="inline-flex rounded-md border border-border p-0.5">
        <Button
          :variant="mode === 'relation' ? 'default' : 'ghost'"
          size="sm"
          class="rounded-sm h-7"
          @click="mode = 'relation'"
        >
          By Relation
        </Button>
        <Button
          :variant="mode === 'entity' ? 'default' : 'ghost'"
          size="sm"
          class="rounded-sm h-7"
          @click="mode = 'entity'"
        >
          By Entity
        </Button>
      </div>

      <Input
        v-model="query"
        placeholder="Search entities or description…"
        class="ml-auto h-8 w-64"
      />
    </div>

    <!-- Content -->
    <ScrollArea class="flex-1 px-6 pb-6">
      <!-- Empty state -->
      <div
        v-if="linkStore.allLinks.length === 0 && !linkStore.loading"
        class="flex flex-col items-center justify-center gap-3 py-20 text-muted-foreground"
      >
        <LinkIcon class="h-10 w-10" />
        <p class="text-sm">No relations yet. Add relations from an entity's detail page.</p>
      </div>

      <!-- By Relation -->
      <Table v-else-if="mode === 'relation'">
        <TableHeader>
          <TableRow>
            <TableHead class="cursor-pointer select-none" @click="toggleSort('from')">
              <span class="inline-flex items-center gap-1">
                From
                <component
                  :is="sortDir === 'asc' ? ChevronUp : ChevronDown"
                  v-if="sortKey === 'from'"
                  class="h-3.5 w-3.5"
                />
              </span>
            </TableHead>
            <TableHead class="cursor-pointer select-none" @click="toggleSort('type')">
              <span class="inline-flex items-center gap-1">
                Type
                <component
                  :is="sortDir === 'asc' ? ChevronUp : ChevronDown"
                  v-if="sortKey === 'type'"
                  class="h-3.5 w-3.5"
                />
              </span>
            </TableHead>
            <TableHead class="cursor-pointer select-none" @click="toggleSort('to')">
              <span class="inline-flex items-center gap-1">
                To
                <component
                  :is="sortDir === 'asc' ? ChevronUp : ChevronDown"
                  v-if="sortKey === 'to'"
                  class="h-3.5 w-3.5"
                />
              </span>
            </TableHead>
            <TableHead>Description</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="row in relationRows"
            :key="row.link.id"
            class="cursor-pointer"
            @click="openRelation(row.link)"
          >
            <TableCell class="font-medium">{{ row.fromName }}</TableCell>
            <TableCell class="text-muted-foreground">{{ row.link.type || '—' }}</TableCell>
            <TableCell class="font-medium">{{ row.toName }}</TableCell>
            <TableCell class="text-muted-foreground max-w-xs truncate">
              {{ row.link.description || '—' }}
            </TableCell>
          </TableRow>
          <TableRow v-if="relationRows.length === 0">
            <TableCell colspan="4" class="text-center text-muted-foreground py-8">
              No relations match your search.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- By Entity -->
      <Table v-else>
        <TableHeader>
          <TableRow>
            <TableHead>Entity</TableHead>
            <TableHead>Type</TableHead>
            <TableHead class="text-right">Relations</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-for="group in entityGroups" :key="group.entity.id">
            <TableRow class="cursor-pointer" @click="toggleExpanded(group.entity.id)">
              <TableCell class="font-medium">
                <span class="inline-flex items-center gap-1.5">
                  <component
                    :is="expanded.has(group.entity.id) ? ChevronDown : ChevronRight"
                    class="h-3.5 w-3.5 text-muted-foreground"
                  />
                  {{ group.entity.name }}
                </span>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ group.entity.type || '—' }}</TableCell>
              <TableCell class="text-right text-muted-foreground">{{ group.links.length }}</TableCell>
            </TableRow>
            <TableRow
              v-for="link in (expanded.has(group.entity.id) ? group.links : [])"
              :key="group.entity.id + ':' + link.id"
              class="cursor-pointer bg-muted/30"
              @click="openRelation(link)"
            >
              <TableCell class="pl-9 font-medium">{{ otherEntityName(link, group.entity.id) }}</TableCell>
              <TableCell class="text-muted-foreground">{{ link.type || '—' }}</TableCell>
              <TableCell class="text-right text-muted-foreground max-w-xs truncate">
                {{ link.description || '—' }}
              </TableCell>
            </TableRow>
          </template>
          <TableRow v-if="entityGroups.length === 0">
            <TableCell colspan="3" class="text-center text-muted-foreground py-8">
              No entities with relations match your search.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </ScrollArea>

  </div>
</template>
