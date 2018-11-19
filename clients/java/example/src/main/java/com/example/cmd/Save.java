package com.example.cmd;

import com.example.ElementProtos;
import com.google.gson.JsonSyntaxException;
import com.google.protobuf.util.JsonFormat;

public class Save extends Client {

    public Save(Options options) {
        super(options);
    }

    public void save() {
        try {
            ElementProtos.Element.Builder builder = ElementProtos.Element.newBuilder();
            JsonFormat.parser().ignoringUnknownFields().merge(options.data, builder);
            ElementProtos.Element ele = newElementService().persistElement(builder.build());
            System.out.println("Element saved:");
            print(ele);
        } catch (JsonSyntaxException e) {
            System.out.println("Invalid json element data " + options.data + " error " + e.getLocalizedMessage());
            System.exit(1);
        } catch (Exception e) {
            System.out.println("Unable to save element error " + e.getLocalizedMessage());
            System.exit(1);
        }
    }

}
