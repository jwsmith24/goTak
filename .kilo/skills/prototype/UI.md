# UI prototype

Create three materially different variants by default, never more than five. Prefer mounting them in an existing page with its real read-only context. Use a clearly named throwaway route only when no suitable host exists.

1. State the design question and variant count.
2. Make variants disagree about layout, information hierarchy, or primary affordance—not merely color or copy.
3. Follow the existing framework, component library, routing, and styling conventions. Add no new library without approval.
4. Select variants through a shareable `?variant=` parameter and a visually separate development-only switcher. Preserve existing data fetching; stub mutations.
5. Report run commands and URLs. Running servers or opening a browser requires separate authorization.
6. After selection, preview production-quality implementation and removal of all switcher/losing-variant code.

Ensure prototype controls cannot appear in production builds, but do not rely on that guard instead of completing cleanup.