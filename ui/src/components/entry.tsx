import React from "react";
import {RouteComponentProps} from "react-router";
import {Query} from "react-apollo";
import gql from "graphql-tag";

import {EntryDetail, entryDetailFragment} from "./entry_detail";
import {GetEntry, GetEntryVariables} from "./__generated__/GetEntry";

const query = gql`query GetEntry($entryId: ID!) {
  getEntry(entryId: $entryId) {
    id
    ...EntryDetailFragment
  }
}
${entryDetailFragment}
`;

interface RouteProps {
  entryId: string
}

export const Entry: React.StatelessComponent<RouteComponentProps<RouteProps>> = ({ match }) => (
  <div className="Entry">
    <Query<GetEntry, GetEntryVariables> query={query} variables={{ entryId: match.params.entryId }}>
      {result => {
        if (result.error) {
          return <p className="error">Error: {result.error.message}</p>
        }
        if (result.loading) {
          return <p className="loading">Loading</p>
        }
        const {data} = result;
        return <EntryDetail entry={data!.getEntry} />;
      }}
    </Query>
  </div>
);