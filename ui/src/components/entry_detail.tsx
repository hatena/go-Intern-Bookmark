import React from "react";
import gql from "graphql-tag";

import {EntryBookmarkFragment} from "./__generated__/EntryBookmarkFragment";
import {EntryDetailFragment} from "./__generated__/EntryDetailFragment";

const entryBookmarkFragment = gql`fragment EntryBookmarkFragment on Bookmark {
  user {
    name
  }
  comment
}
`;

interface EntryBookmarkProps {
  bookmark: EntryBookmarkFragment
}
const EntryBookmark: React.StatelessComponent<EntryBookmarkProps> = ({ bookmark }) => (
  <div className="EntryBookmark">
    <h2>{bookmark.user.name}</h2>
    <p>{bookmark.comment}</p>
  </div>
);


export const entryDetailFragment = gql`fragment EntryDetailFragment on Entry {
  title
  url
  bookmarks {
    id
    ...EntryBookmarkFragment
  }
}
${entryBookmarkFragment}
`;

interface EntryDetailProps {
  entry: EntryDetailFragment
}
export const EntryDetail: React.StatelessComponent<EntryDetailProps> = ({ entry }) => (
  <div className="EntryDetail">
    <h1>{entry.title}</h1>
    <p><a href={entry.url}>{entry.url}</a></p>
    <ul>
      {entry.bookmarks.map(bookmark => (<li key={bookmark.id}>
        <EntryBookmark bookmark={bookmark} />
      </li>))}
    </ul>
  </div>
);
