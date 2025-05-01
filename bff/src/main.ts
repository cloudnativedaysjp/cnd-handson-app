import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Transport, MicroserviceOptions } from '@nestjs/microservices';
import * as path from 'path';
async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule,
    {
      options: {
        package: 'project',
        protoPath: path.join(__dirname, 'project/project.proto'),
      },
    },
  );
  await app.listen();
}
bootstrap();
