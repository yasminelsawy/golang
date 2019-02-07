package com.example.yasmi.unirest;

import android.content.Context;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.TextView;

import com.github.library.bubbleview.BubbleTextView;

import java.util.ArrayList;

/**
 * Created by yasmi on 11/29/2017.
 *  Adapter. In android, an adapter is a bridge between UI component and data source that helps us to fill data in the UI component.
 *  It holds the data and send the data to adapter view then view can takes the data from the adapter view and shows the data on
 *  different views like as list view, grid view, spinner etc.
 */

public class ChatAdapter extends BaseAdapter {
    ArrayList<String> serverMessage=new ArrayList<>();
    ArrayList<String> UserMessage=new ArrayList<>();

    //Instantiates a layout XML file into its corresponding View objects.
    //Law 3ndi Layout its connected bel inflator and the inflator is connected bel CODE
    LayoutInflater mInflaer;
    public ChatAdapter(Context context, ArrayList<String> ServerMessage, ArrayList<String> userMessage ) {
        if(context==null){
            Log.e("CONTEXT","NULLL" );
        }
        for(int i=0;i<ServerMessage.size();i++){
            serverMessage.add(ServerMessage.get(i));

        }
        for (int i=0;i<userMessage.size();i++){
            UserMessage.add(userMessage.get(i));
        }

        mInflaer=(LayoutInflater) context.getSystemService(Context.LAYOUT_INFLATER_SERVICE);


    }
    @Override
    public int getCount() {
        return serverMessage.size() ;
    }

    @Override
    public Object getItem(int i) {
        return serverMessage.get(i);
    }

    @Override
    public long getItemId(int i) {
        return i;
    }

    @Override
    /*This function is automatically called when the list item view is ready to be displayed or about to be displayed.
     In this function we set the layout for list items using LayoutInflater class and then add the data to the views
    * */
    public View getView(int i, View view, ViewGroup viewGroup) {
        View v=mInflaer.inflate(R.layout.chatmessagelistview , null);
       final TextView server =(BubbleTextView) v.findViewById(R.id.wlcmMessagee);
        TextView user =(BubbleTextView) v.findViewById(R.id.userMessagee);
      final String serverMsg =serverMessage.get(i);
        if(serverMsg.contains("https")){

                    String href = serverMsg.substring(44, serverMsg.length());

                    server.setText(href);
                    Log.e("IT HTTPS", "link");

        }else{
            server.setText(serverMsg);
        }

        Log.i("GET VIEW STRING server", serverMsg);
        if(i<=UserMessage.size()-1){
            String userMsg =UserMessage.get(i);
            Log.i("GET VIEW STRING User", userMsg);
            user.setText(userMsg);
        }else{
            user.setVisibility(View.INVISIBLE);
            TextView textView=(TextView)  v.findViewById(R.id.textView);
            textView.setVisibility(View.INVISIBLE);
        }


        return v;
    }
}
