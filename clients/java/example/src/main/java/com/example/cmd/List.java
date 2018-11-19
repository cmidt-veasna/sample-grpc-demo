package com.example.cmd;

import com.example.ElementProtos;
import com.google.gson.JsonSyntaxException;
import com.google.protobuf.util.JsonFormat;

public class List extends Client {

    public List(Options options) {
        super(options);
    }

    public void list() {
        try {
            ElementProtos.ElementFilter.Builder builder = ElementProtos.ElementFilter.newBuilder();
            JsonFormat.parser().ignoringUnknownFields().merge(options.filter, builder);
            ElementProtos.Elements eles = newElementService().listElement(builder.build());
            if (eles == null || eles.getElementsList() == null || eles.getElementsCount() == 0) {
                System.out.println("No element found.");
                System.exit(0);
            }
            print(eles);
        }catch (JsonSyntaxException e) {
            System.out.println("Invalid json element filter " + options.filter + " error " + e.getLocalizedMessage());
            System.exit(1);
        }catch (Exception e) {
            System.out.println("Unable to list element error " + e.getLocalizedMessage());
            System.exit(1);
        }
    }
}
