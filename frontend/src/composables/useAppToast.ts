import { toast } from 'vue-sonner'

/**
 * Thin wrapper around vue-sonner that enforces a consistent two-line shape:
 *
 *   ┌─────────────────────────────────────────┐
 *   │ ✕  Failed to create world               │  ← title  (what the user was doing)
 *   │    folder is not empty: C:\My\Folder    │  ← detail (raw Go error — selectable)
 *   └─────────────────────────────────────────┘
 *
 * The detail line is plain selectable text so users can copy-paste it when
 * reporting bugs.
 */
export function useAppToast() {
  /**
   * Show a red error toast.
   * @param title  Short human label: "Failed to create world"
   * @param detail Raw error string from Go (displayed smaller, fully selectable)
   */
  function errorToast(title: string, detail?: string) {
    toast.error(title, { description: detail })
  }

  /**
   * Show a green success toast.
   * @param title  Short human label: "World saved"
   * @param detail Optional supplementary text
   */
  function successToast(title: string, detail?: string) {
    toast.success(title, { description: detail })
  }

  return { errorToast, successToast }
}
