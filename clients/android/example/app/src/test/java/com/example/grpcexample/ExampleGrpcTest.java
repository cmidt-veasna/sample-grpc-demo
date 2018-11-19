package com.example.grpcexample;

import com.example.ElementProtos;
import com.example.ElementServiceGrpc;

import org.junit.Test;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotEquals;

/**
 * Example local unit test, which will execute on the development machine (host).
 *
 * @see <a href="http://d.android.com/tools/testing">Testing documentation</a>
 */
public class ExampleGrpcTest {
    @Test
    public void persistElement() {
        ManagedChannel mc = ManagedChannelBuilder.forAddress("localhost", 8080).usePlaintext().build();
        ElementServiceGrpc.ElementServiceBlockingStub blockStub = ElementServiceGrpc.newBlockingStub(mc);

        ElementProtos.Element element = ElementProtos.Element.newBuilder()
                .setName("Sample")
                .setAge(24)
                .setStatus(11)
                .build();

        ElementProtos.Element elementN = blockStub.persistElement(element);
        assertEquals(element.getName(), elementN.getName());
        assertEquals(element.getAge(), elementN.getAge());
        assertEquals(element.getStatus(), elementN.getStatus());
        assertNotEquals("", elementN.getId());
        assertNotEquals("", elementN.getCreatedAt());
        assertNotEquals("", elementN.getUpdatedAt());
    }

    @Test
    public void listElement() {
        ManagedChannel mc = ManagedChannelBuilder.forAddress("localhost", 8080).usePlaintext().build();
        ElementServiceGrpc.ElementServiceBlockingStub blockStub = ElementServiceGrpc.newBlockingStub(mc);

        ElementProtos.ElementFilter ef = ElementProtos.ElementFilter.newBuilder()
                .setAge("[12, 40]")
                .setStatus("{11, 13, 15}")
                .build();

        ElementProtos.Elements elements = blockStub.listElement(ef);
        assertNotEquals(null, elements);
    }
}
