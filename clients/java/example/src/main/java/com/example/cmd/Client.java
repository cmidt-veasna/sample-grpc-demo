package com.example.cmd;

import com.example.ElementProtos;
import com.example.ElementServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

class Client {

    Options options;

    Client(Options options) {
        this.options = options;
    }

    ElementServiceGrpc.ElementServiceBlockingStub newElementService() {
        ManagedChannel channel = ManagedChannelBuilder.forAddress(options.address, options.port).usePlaintext().build();
        return ElementServiceGrpc.newBlockingStub(channel);
    }

    void print(ElementProtos.Element e) {
        System.out.printf("Id:\t\t%s\n", e.getId());
        System.out.printf("Name:\t\t%s\n", e.getName());
        System.out.printf("Age:\t\t%d\n", e.getAge());
        System.out.printf("Status:\t\t%d\n", e.getStatus());
        System.out.printf("CreateAt:\t%s\n", e.getCreatedAt());
        System.out.printf("UpdateAt:\t%s\n", e.getUpdatedAt());
    }

    void print(ElementProtos.Elements es) {
        for (ElementProtos.Element e : es.getElementsList()) {
            print(e);
            System.out.println("-------");
        }
    }

}
