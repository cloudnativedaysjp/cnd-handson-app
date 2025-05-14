import { ApolloServer } from "apollo-server"
import { loadSchemaSync } from "@graphql-tools/load";
import { GraphQLFileLoader } from "@graphql-tools/graphql-file-loader";

const schema = loadSchemaSync("./schema.graphql", {
  loaders: [new GraphQLFileLoader()],
});

// サーバーを起動する
const server = new ApolloServer({ schema, mocks: true });

server.listen().then(({ url }) => {
  console.log(`🚀 Server ready at ${url}`);
});