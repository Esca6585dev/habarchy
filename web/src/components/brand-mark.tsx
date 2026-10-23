// BrandMark is the Habarçy SMS logo glyph (concept #3: a message bubble with a
// delivered double-check). It draws in currentColor, so a parent sets the
// colour — e.g. a blue tile with text-primary-foreground. One source of truth
// for the mark across the sidebar, login and anywhere else it appears.
export function BrandMark({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      className={className}
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M21 14a2 2 0 0 1-2 2H9l-4 3.5V5a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />
      <path d="M7.5 10.2l2.4 2.4 4-4.2" />
      <path d="M12.3 12.6l.3.3 4-4.2" />
    </svg>
  );
}
