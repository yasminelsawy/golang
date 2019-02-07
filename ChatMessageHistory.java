package com.example.yasmi.unirest;

import java.util.ArrayList;

public class ChatMessageHistory {
    ArrayList<String> serverMessage=new ArrayList<>();
    ArrayList<String> UserMessage=new ArrayList<>();

    public ChatMessageHistory(){

        super();

    }

    public void setMessageServerHistory(String servermessage) {
        serverMessage.add(servermessage);

    }

    public int getSizeOfServerHistory(){
        return serverMessage.size();
    }


    public int getSizeOfUserHistory(){
        return UserMessage.size();
    }

    public void setMessageUserHistory(String Usermessage) {
        UserMessage.add(Usermessage);
    }

    public ArrayList<String> getServerMessage() {

        return serverMessage;
    }
    public void reset(){
        ArrayList<String> serverMessage=new ArrayList<>();
        ArrayList<String> UserMessage=new ArrayList<>();
    }

    public ArrayList<String> getUserMessage() {

        return UserMessage;
    }

}
