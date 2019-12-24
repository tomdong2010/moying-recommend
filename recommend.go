if(cmd == "new"){
        count = 0

        my_user_upvideo_lock.RLock()
        all_len := uint32(len(my_user_upvideo))
        my_user_upvideo_lock.RUnlock()
        for short_id := uint32(0) ;short_id<all_len; short_id++ {
            my_user_upvideo_lock.RLock()
            vid_1 := my_user_upvideo[short_id].vid
            var my_user_upvideo_t = my_user_upvideo[short_id]
            my_user_upvideo_lock.RUnlock()

            vid_watch_list.RLock()
            _, isset := vid_watch_list.Items[vid_1]
            vid_watch_list.RUnlock()
            
            if(debugmode ==2){
               //fmt.Println("count:",count,",short_id:",short_id,"vid",vid_1,my_user_upvideo_t.limit_play_times,my_user_upvideo_t.good_number)
            }

            if(!isset && vid_1>0){

                if ( my_user_upvideo_t.limit_play_times <=0){
                    continue
                }

                if(my_user_upvideo_t.limit_play_times > 0){
                    my_user_upvideo_t.limit_play_times = my_user_upvideo_t.limit_play_times- 1
                    my_user_upvideo_lock.Lock()
                    my_user_upvideo[short_id] = my_user_upvideo_t
                    my_user_upvideo_lock.Unlock()
                }

                if(int64(count) <= limit1 && limit1 !=0){
                    count ++
                    continue
                }

                count ++
                search_result += strconv.FormatInt( int64(vid_1) , 10) + ","
                
                if(debugmode > 0){
                    fmt.Println("count:",count,",videoo.id:",vid_1)
                }
                // ,search_result
                if int64(count) >= limit1 + limit2{
                    break
                }
                        
            }
        }
    }else{

        //fmt.Println("My_recommend_data len: ",len(My_recommend_data))
        fmt.Println("find_ids  len: ",len(find_ids))
        //Openid_like_list_lock.RLock()
        like_list := map[uint32]bool {}
        like_list_recent := map[uint32]bool {}
        like_list0 ,ok := Openid_like_list.Get(openid) //Openid_like_list[openid]
        if ok {
            like_listt := like_list0.([]uint32)
            for idx,v := range like_listt {
                like_list[v] = true
                if idx > len(like_listt)- 20{
                    like_list_recent[v]= true
                }
            }
        }
        Company_admin_like_list_t := map[uint32]uint8 {}
        company_admin_openid_lock.RLock()
        company_admin_openid_1 , ok := company_admin_openid[string(openid)]
        company_admin_openid_lock.RUnlock()
        if ok {

            admin_openids := company_admin_openid_1.openid
            for k := 0; k < len(admin_openids); k++ {
                openid_t := admin_openids[k]
                like_list0 ,ok := Openid_like_list.Get(openid_t) // Openid_like_list[openid_t]
                if ok {
                    for _,v := range like_list0.([]uint32) {
                        t,ok2:= Company_admin_like_list_t[v]
                        tt := uint8(1)
                        if ok2{
                            tt = t
                        }
                        if tt <2{
                            Company_admin_like_list_t[v] = tt + 1
                        }
                    }
                }
            }
        }
        
       
        //Openid_like_list_lock.RUnlock()
        for i ,id1 := range find_ids {
            
            watch_time := vid_watch_list2[uint32(id1)] & 0x00ff
            //
            var score0 float64

            score0 = 0
            if watch_time >120 {
                score0 = 2.79 + float64(watch_time - 120)*0.067/8
            }else if watch_time >58 {
                score0 = 2 + float64(watch_time - 58)*0.067/6
            }else if watch_time >25 {
                score0 = 1.375 + float64(watch_time - 25)*0.067/4
            }else if watch_time> 14 {
                score0 = 1.0+float64(watch_time - 14)*0.067/2
            }else if watch_time> 7 {
                score0 = 0.5+float64(watch_time - 7)*0.067
            }else if watch_time> 1 {
                score0 = float64(watch_time )*0.01
            }else {
                //continue
            }
            /*else if watch_time >2{
                score0 = float64(watch_time )*0.05
            }*/
            if (debugmode > 0) {
                //fmt.Print("vid:",id1,",watch_time:",watch_time,",score0:",score0,"\n")
            }

            fav_vid, isset := like_list[  id1  ]
            //Openid_like_list_lock.RUnlock()
            _ = fav_vid

            _,isrecent := like_list_recent[id1]
            if isrecent {
                score0 += 4.0
            }
                    
            couter1,isset_t := Company_admin_like_list_t[id1]

            if isset_t {
                score0 += 2.0 * float64(couter1)
                //fmt.Println("find company_admin_openid openid:",openid," openid_t:",admin_openids[k],strconv.Itoa(int(id1)))
            }
            //fmt.Println(string(openid)+":"+strconv.Itoa(int(id1)),",fav_vid:",fav_vid,",isset:",isset)
            if (isset){
                score0 += 2.0
                if (debugmode > 0) {
                    //fmt.Print("openid:",openid,",vid:",id1,", liked score+2.0 \n")
                }
            }

            if (score0>10) {
                score0 = 10
            }

            r_idx := uint32(id1)
            My_recommend_data_lock.RLock()
            my_recommend_t , ok := My_recommend_data[r_idx]
            My_recommend_data_lock.RUnlock()
            if(!ok){
                continue
            }
            vid_index := [] int {2,4,3,0,1}
            sum_score := float64(0)
            
            for j:=0;j<5;j++{
                vid_recommend := my_recommend_t.vid[vid_index[j]]
                vid_watch_list.RLock()
                tm, isset := vid_watch_list.Items[vid_recommend]
                _,ishot := Super_hot_video[vid_recommend]
                if ishot {
                    continue
                }
                vid_watch_list.RUnlock()
                if (debugmode > 0) {
                //fmt.Print("vid:",video_id_0[find_ids[i]],",tm:",tm,",isset:",isset,"\n")
                }
                if( !isset ){
                     for k:=0;k<5;k++{
                        sum_score += float64(my_recommend_t.score[k])
                   }
                    if (debugmode > 2){
                        fmt.Println("vid:",id1,",find_ids[i]",find_ids[i],"video_id_0[find_ids[i]]",vid_recommend)
                    }
                    if(vid_recommend != 0){
                        if (debugmode > 2){
                            fmt.Print("vid:",id1,",video_id_0[find_ids[i]]:",vid_recommend,"\n")
                        }
                        score1 := score0 * float64(my_recommend_t.score[vid_index[j]]) /sum_score +  float64(my_recommend_t.hot_score-730) / 600 //*(1.+float64(0.05)*float64(5-j))

                        video0 := video { 
                            id: vid_recommend,
                            score : int32(score1*1000)}
                        video_select = append(video_select, video0)
                        break
                    }
                }else if(debugmode >0){
                    _ = tm
                }
            }
            if len(video_select) >10000 {
                break
            }
        }
        

        fmt.Println("video_select short order len: ",len(video_select))
        sort.Stable(video_select)

        count = 0

        for _,videoo := range video_select {
            vid_watch_list.RLock()
            t ,isset := vid_watch_list.Items[videoo.id] 
            vid_watch_list.RUnlock()
            _=t
            if (! (videoo.id >0) || isset ){
                continue
            }

            if(int64(count) <= limit1 && limit1 !=0){
                count ++
                continue
            }

            count ++
            search_result += strconv.FormatInt( int64(videoo.id) , 10) + ","
            
            if(debugmode > 0){
                fmt.Println("count:",count,",videoo.id:",videoo.id,",videoo.score:",videoo.score,",isset:",isset)
            }
            // ,search_result
            if int64(count) >= limit1 + limit2{
                break
            }
        }
    }