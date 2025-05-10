import { Module } from "@nestjs/common";
import { AppService } from "./app.service";
import { ClientsModule, Transport } from "@nestjs/microservices";
import { USER_PACKAGE_NAME } from "./user/user";
import { SESSION_PACKAGE_NAME } from "./session/session";
import * as path from "path";
import { AppResolver } from "./app.resolver";
import { AppController } from "./app.controller";
@Module({
  imports: [
    ClientsModule.register([
      {
        name: "USER_PACKAGE",
        transport: Transport.GRPC,
        options: {
          url: "localhost:5000",
          package: USER_PACKAGE_NAME,
          protoPath: path.join(__dirname, "user/user.proto"),
        },
      },
    ]),
  ],
  controllers: [AppController],
  providers: [AppService, AppResolver],
})
export class AppModule {}