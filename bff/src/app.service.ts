import { Inject, Injectable, OnModuleInit } from "@nestjs/common";
import { ClientGrpc } from "@nestjs/microservices";
import { Observable, lastValueFrom } from "rxjs";
import { adjustRpcResponse } from "./utils/convertObservableToPromise";

@Injectable()
export class AppService implements OnModuleInit {

}