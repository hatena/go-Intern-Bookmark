import React from "react";
import {Link} from "react-router-dom";
import gql from "graphql-tag";

import {EntryListItemFragment} from "./__generated__/EntryListItemFragment";
import {EntryListFragment} from "./__generated__/EntryListFragment";

const entryListItemFragment = gql`fragment EntryListItemFragment on Entry {
  id
  title
  url
}`;

interface EntryListItemProps {
  entry: EntryListItemFragment
}
export const EntryListItem: React.StatelessComponent<EntryListItemProps> = ({ entry }) => (
  <div className="EntryListItem">
    <Link to={`/entry/${entry.id}`}>{entry.title}</Link>
    <span> - </span>
    <a href={entry.url}>{entry.url}</a>
  </div>
);

export const entryListFragment = gql`fragment EntryListFragment on Entry {
  id
  ...EntryListItemFragment
}
${entryListItemFragment}
`;

interface EntryListProps {
  entries: EntryListFragment[]
}
export const EntryList: React.StatelessComponent<EntryListProps> = ({ entries }) => (
  <ul className="EntryList">
    {entries.map(entry => (<li key={entry.id}><EntryListItem entry={entry} /></li>))}
  </ul>
);
