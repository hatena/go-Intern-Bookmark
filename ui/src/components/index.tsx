import React from "react";
import {Query} from "react-apollo";
import gql from "graphql-tag";

import {EntryList, entryListFragment} from "./entries";
import {ListEntries} from "./__generated__/ListEntries";

const query = gql`query ListEntries {
  listEntries {
    ...EntryListFragment
  }
}
${entryListFragment}
`;

export const Index: React.StatelessComponent = () => (
  <div className="Index">
    <h1>Entries</h1>
    <Query<ListEntries> query={query}>
      {result => {
        if (result.error) {
          return <p className="error">Error: {result.error.message}</p>
        }
        if (result.loading) {
          return <p className="loading">Loading</p>
        }
        const { data } = result;
        return <EntryList entries={data!.listEntries} />;
      }}
    </Query>
  </div>
);