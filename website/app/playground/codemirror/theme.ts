import { githubDarkInit } from '@uiw/codemirror-theme-github';
import { tags } from '@lezer/highlight';

// Keep example previews and the editable playground on the same palette.
export const darijaScriptTheme = githubDarkInit({
  styles: [
    { tag: tags.keyword, color: '#b972cf' },
    {
      tag: [tags.function(tags.variableName), tags.function(tags.propertyName)],
      color: '#e5c07b',
    },
    { tag: [tags.bool, tags.null], color: '#ffab70' },
    { tag: tags.operator, color: '#c9d1d9' },
    { tag: tags.string, color: '#93bd76' },
    { tag: tags.escape, color: '#79c0ff' },
    { tag: tags.number, color: '#b1855d' },
  ],
});
